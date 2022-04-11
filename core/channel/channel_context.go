package channel

import (
	"camellia/core/datapack"
	"sync"
)

type ConnContext struct {
	isInit   bool
	initLock sync.Mutex

	WriteChan  chan<- []byte

	//handler chain
	Head, Tail HandlerContext
}

//InitHandlerContext init default and handlerContext Initializer provider func
func (ctx *ConnContext) InitHandlerContext(providers ...Initializer) {
	if ctx.isInit {
		return
	}
	ctx.initLock.Lock()
	defer ctx.initLock.Unlock()

	//double check
	if ctx.isInit {
		return
	}

	ctx.Head = HandlerContext{handler: &HeadDataHandler{}}
	ctx.Tail = HandlerContext{handler: &TailDataHandler{}}

	ctx.Head.next = &ctx.Tail
	ctx.Tail.pre = &ctx.Head

	//add other handler context
	for _, provider := range providers {
		ctx.AddHandlerContext(provider())
	}

	ctx.isInit = true
}

//AddHandler add handler to last(before tail)
func (ctx *ConnContext) AddHandler(handler DataHandler) {
	ctx.AddHandlerContext(HandlerContext{
		handler: handler,
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
	handler   DataHandler
	pre, next *HandlerContext
}

func (h *HandlerContext) Fire(ctx *ConnContext, pkg datapack.Message) {
	h.handler.Exec(ctx, pkg)
	if h.next != nil {
		h.next.Fire(ctx, pkg)
	}
}


//Initializer handlerContext provider
type Initializer func() HandlerContext

func DispatchHandlerContextFunc() HandlerContext {
	//todo impl DispatchHandler replace StdDataHandler
	return HandlerContext{handler: &StdDataHandler{}}
}

func AckHandlerContextFunc() HandlerContext {
	//todo impl AckHandler replace StdDataHandler
	return HandlerContext{handler: &StdDataHandler{}}
}


