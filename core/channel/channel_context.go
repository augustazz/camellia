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

func (ctx *ConnContext) InitHandlerContext() {
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

	ctx.isInit = true

	ctx.AddHandler(&StdDataHandler{})
}

//AddHandler add handler to last(before tail)
func (ctx *ConnContext) AddHandler(handler DataHandler) {
	ctx.AddHandlerContext(&HandlerContext{
		handler: handler,
	})
}

//AddHandlerContext add ctx to last(before tail)
func (ctx *ConnContext) AddHandlerContext(handler *HandlerContext) {
	if !ctx.isInit {
		ctx.InitHandlerContext()
	}

	tmp := ctx.Tail.pre
	tmp.next = handler
	handler.pre = tmp
	handler.next = &ctx.Tail
	ctx.Tail.pre = handler
}

//HandlerContext wrap handlers with linklist
type HandlerContext struct {
	handler   DataHandler
	pre, next *HandlerContext
}

func (h *HandlerContext) Fire(ctx *ConnContext, pkg *datapack.TcpPackage) {
	h.handler.Exec(ctx, pkg)
	if h.next != nil {
		h.next.Fire(ctx, pkg)
	}
}

type HandlerContextInitializer struct {

}

func (i *HandlerContextInitializer) Initializer(ctx ...HandlerContext) {

}
