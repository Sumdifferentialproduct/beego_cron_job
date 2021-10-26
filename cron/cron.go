package cron
//封装任务结构体
import (
	"github.com/gorhill/cronexpr"
	"sort"
	"time"
)

//为了程序的可扩展性
//每个job运行方式是不一样的。所以定义接口类型
//通过重写job接口run()，实现不同的job

type Job  interface {
	//定义job的执行方法
	Run()
}

type Entry struct {
	//任务时间表达式
	Schedule    cronexpr.Expression
	//下一次执行时间
	Next  time.Time
	//上一次执行时间
	Prev  time.Time
	Job
}

//定义类型，用域管道删除任务，返回布尔值代表是否要删除
//true  删除    ，    false   不删除
type RemoveCheckFunc  func(e  *Entry)  bool

//对照图继续封装
type  Cron struct {
	//任务列表
	entries   []*Entry
	//需要一些管道，处理任务调度相关的功能
	//添加管道
	add  chan	*Entry
	//删除管道
	remove   chan  RemoveCheckFunc
	//复制管道
	snapshot  chan []*Entry
	//任务是否在运行的标识
	running   bool
}

//对于需要make的，比如chan，提供一下初始化方法，便于使用
//定义构造方法，也就是初始化值


//启动Cron

func (c *Cron)  Start  (){
	c.running =true
	//启动任务后，应该又一个调度协程，不断去调度，去处理整个的任务执行逻辑
	go  c.run()
}

//调度协程
//1.遍历所有的任务，根据当前时间，计算下一次的执行时间
//2.循环下一次执行时间，对所有任务根据next进行排序
//3.调度协程根据不同管道的数据做不同的调度
func (c *Cron)  run  (){
	//获取当前时区的时间
	now :=  time.Now().Local()
	//遍历所有任务，得到所有任务的下一次执行时间
	for _,entry := range  c.entries{
		entry.Next = entry.Schedule.Next(now)
	}
	for {
		//根据entry的next时间进行排序
		sort.Sort(byTime(c.entries))
		//定义最近一次要执行的时间
		var  nearTime  time.Time
		//获取最近一次将要执行的时间
		if len(c.entries)==0 || c.entries[0].Next.IsZero(){
			//把最近一次要执行的时间加10年，相当于休眠
			nearTime = now.AddDate(10,0,0)
		}else{
			nearTime = c.entries[0].Next
		}
		//调度协程做不同的调度  最近要执行的任务到时间了自动执行，根据管道进行调度
		select {
		case now = <-time.After(nearTime.Sub(now)):
			//遍历所有的任务，同一时间可能又多个任务执行
			for _ ,e := range  c.entries{
				if e.Next  !=nearTime{
					break
				}
				//协程执行任务，并发调度
				go e.Job.Run()
				//xxx任务。 next =9.51  prev =10.01
				//xxx任务。 next =10.01  prev =9.51
				e.Prev = e.Next
				e.Next = e.Schedule.Next(nearTime)
			}
			//监控任务添加，相当于开启任务
			case newEntry := <-c.add:
			c.entries = append(c.entries,newEntry)
			//计算任务时间
			newEntry.Next = newEntry.Schedule.Next(now)
			//监控删除任务的管道，相当于停止任务
			//type RemoveCheckFunc  func(e  *Entry)  bool
			case cb := <-c.remove:
			//创建一个空切片
			newEntries := make([]*Entry,0)
			for _,e    :=range c.entries{
				//cb(e)返回true，代表遍历到了要删除的任务
				//cb(e)返回false，代表不删除任务，取反为true，则进入if判断
				if !cb(e){
					newEntries =append(newEntries,e)
				}
			}
			c.entries = newEntries
			//监听复制管道
		case <-c.snapshot:
			c.snapshot	<-c.entrySnapshot()

		}
		//更新下当前时间
		now = time.Now().Local()
	}

}

//复制任务列表
func (c *Cron)entrySnapshot()[]*Entry{
	entries := []*Entry{}
	for _,e := range c.entries {
		entries = append(entries,&Entry{
			Schedule: e.Schedule,
			Next :e.Next,
			Prev: e.Prev,
			Job :	e.Job,
		})
	}
	return entries
}

//实现sort排序方法
//需要一个entry任务切片
type   byTime  []*Entry
func (s byTime) Len()int{
	return len(s)
}
func (s byTime)Swap(i,j int){
	s[i],s[j]=s[j],s[i]

}
//按照任务的下次执行时间排序。最近要执行的放在前面，小的放在前面
//若i的数据小于j的数据，返回true，升序排序,不交换
func (s byTime)Less(i,j int)bool{
	//IsZero()表示Time类型的零值00：00：00，代表下一次执行时间已过，或者时间无效
	if s[i].Next.IsZero(){
		return false
	}
	if s[j].Next.IsZero(){
		return true
	}
	//比较i时间，是否在j时间之前，不在返回false，则进行交换
	return s[i].Next.Before(s[j].Next)
}


//还需要封装外部调用的方法
//任务添加
//任务是啥,任务多长时间执行一次
func (c *Cron)AddJob(spec  string,job  Job) error{
	schedule, err := cronexpr.Parse(spec)
	if err !=nil{
		return  err
	}
	//构建任务对象
	entry := &Entry{
		Schedule :*schedule,
		Job : job,
	}
	//任务添加进任务列表
	if  c.running {
		c.entries =append(c.entries,entry)
	}
	c.add<- entry
	return nil
}

//任务删除
//type RemoveCheckFunc  func(e  *Entry)  bool
func (c *Cron)RemoveJob(cb  RemoveCheckFunc) {
	//添加进管道
	c.remove <-cb

}

//任务复制

func (c *Cron)Entries() []*Entry {
	if c.running{
		c.snapshot <- nil
		entrys := 	<-c.snapshot
		return  entrys
	}
	return c.entrySnapshot()
}
