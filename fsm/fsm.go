/*
有限状态机(Finite-state machine, 简写FSM)又可以称作有限状态自动机。
它必须是可以附着在某种事物上的，且该事物的状态是有限的，通过某些触发事件，会让其状态发生转换。
为此，有限状态机就是描述这些有限的状态和触发事件及转换行为的数学模型。

有限状态机有两个必要的特点，一是离散的，二是有限的。

状态(State)：事物的状态，包括初始状态和所有事件触发后的状态
事件(Event)：触发状态变化或者保持原状态的事件
行为或转换(Action/Transition)：执行状态转换的过程
检测器(Guard)：检测某种状态要转换成另一种状态的条件是否满足

*/

package fsm

import (
	"fmt"
	"sync"
)

type FSMState string            //状态
type FSMEvent string            //事件
type FSMHandler func() FSMState //匿名函数类型,返回新状态

type FSM struct {
	mu       sync.Mutex                           //排它锁
	state    FSMState                             //当前状态
	handlers map[FSMState]map[FSMEvent]FSMHandler //处理集合,每个状态都可以发出有限个事件,执行有限个处理
}

//获取状态
func (f *FSM) getState() FSMState {
	return f.state
}

//设置当前状态
func (f *FSM) setState(newState FSMState) {
	f.state = newState
}

//为状态添加事件处理方法
func (f *FSM) AddHandler(state FSMState, event FSMEvent, handler FSMHandler) *FSM {
	if _, ok := f.handlers[state]; !ok {
		f.handlers[state] = make(map[FSMEvent]FSMHandler)
	}
	if _, ok := f.handlers[state][event]; ok {
		fmt.Printf("[警告] 状态(%s)事件(%s)已定义过", state, event)
	}
	f.handlers[state][event] = handler
	return f
}

//事件处理
func (f *FSM) Call(event FSMEvent) FSMState {
	f.mu.Lock()
	defer f.mu.Unlock()
	events := f.handlers[f.getState()]
	if events == nil {
		return f.getState()
	}
	if fn, ok := events[event]; ok {
		oldstate := f.getState()
		f.setState(fn())
		newState := f.getState()
		fmt.Println("状态从[", oldstate, "]切换到[", newState, "]")
	}
	return f.getState()
}

//实例化FSM
func NewFSM(initState FSMState) *FSM {
	return &FSM{
		state:    initState,
		handlers: make(map[FSMState]map[FSMEvent]FSMHandler),
	}
}
