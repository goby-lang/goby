package vm

func initSpecClass(vm *VM) {
	vm.mainThread.execGobyLib("spec.gb")
}
