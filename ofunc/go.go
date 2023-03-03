/*
 *  * Copyright (c) 2023 guojia99 All rights reserved.
 *  * Created: 2023/3/3 下午6:44.
 *  * Author: guojia(https://github.com/guojia99)
 */

package ofunc

import `sync`

type WaitTool struct {
	wg sync.WaitGroup
}

func (w *WaitTool) Go(fn func()) {
	w.wg.Add(1)
	go func() {
		fn()
		w.wg.Done()
	}()
}

func (w *WaitTool) Wait() {
	w.wg.Wait()
}
