package channel

import (
	"camellia/core/datapack"
	"camellia/core/enums"
	"sync"
	"time"
)

type ConnContext struct {
	isInit   bool
	initLock sync.Mutex

	State     enums.ConnState
	WriteChan chan<- []byte
	RandomStr string

	//handler chain
	Abort      bool //中断传递
	Head, Tail HandlerContext

	ConnectTime    time.Time
	LastReadTime   time.Time
	LastWriteTime  time.Time
}

//InitHandlerContext init default and handlerContext Initializer provider func
func (ctx *ConnContext) InitHandlerContext(handlers ...func(ctx *ConnContext, msg datapack.Message)) {
	if ctx.isInit {
		return
	}
	ctx.initLock.Lock()
	defer ctx.initLock.Unlock()

	//double check
	if ctx.isInit {
		return
	}

	ctx.Head = HandlerContext{Handler: HeadDataHandlerFunc}
	ctx.Tail = HandlerContext{Handler: TailDataHandlerFunc}

	ctx.Head.next = &ctx.Tail
	ctx.Tail.pre = &ctx.Head

	//add other handler context
	for _, handler := range handlers {
		ctx.AddHandler(handler)
	}

	ctx.ConnectTime = time.Now()
	ctx.isInit = true
}

//AddHandler add handler to last(before tail)
func (ctx *ConnContext) AddHandler(handlerFunc func(ctx *ConnContext, msg datapack.Message)) {
	ctx.AddHandlerContext(HandlerContext{
		Handler: handlerFunc,
	})
}

//AddHandlerContext add ctx to last(before tail)
func (ctx *ConnContext) AddHandlerContext(handler HandlerContext) {
	tmp := ctx.Tail.pre
	tmp.next = &handler
	handler.pre = tmp
	handler.next = &ctx.Tail
	ctx.Tail.pre = &handler
}

//func (ctx *ConnContext) HandlerContextInitializer() func(apply ...Initializer) {
//	return func(apply ...Initializer) {
//		for _, h := range apply {
//			ctx.AddHandlerContext(h())
//		}
//	}
//}

//HandlerContext wrap handlers as linklist node
type HandlerContext struct {
	Handler   func(ctx *ConnContext, msg datapack.Message)
	pre, next *HandlerContext
}

func (h *HandlerContext) Fire(ctx *ConnContext, msg datapack.Message) {
	//if abort且不是tai节点， 直接传给下一个，h.next==nil表示tail
	if ctx.Abort && h.next != nil {
		h.next.Fire(ctx, msg)
	}

	//执行func
	h.Handler(ctx, msg)

	//传递
	if h.next != nil {
		h.next.Fire(ctx, msg)
	}
}
