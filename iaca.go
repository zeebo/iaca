package iaca

//go:noinline
func padStart()   {}
func padStart_()  { padStart() }
func padStart__() { padStart_() }

func Start() { padStart__() }

//go:noinline
func padStop()   {}
func padStop_()  { padStop() }
func padStop__() { padStop_() }

func Stop() { padStop__() }
