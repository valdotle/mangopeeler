package internal

type Pool struct {
	addJob          chan string
	removeJob, wait chan any
	jobs            []string
	size            uint
}

func (p Pool) Add(path string) {
	p.addJob <- path
}

func (p Pool) Remove() {
	p.removeJob <- nil
}

func (p Pool) Finish() {
	<-p.wait
}

func NewPool(dirThreads, dirEntryThreads uint, processor func(string)) Pool {
	p := Pool{make(chan string, dirEntryThreads*dirThreads), make(chan any, dirThreads), make(chan any), nil, dirThreads}

	go p.run(processor)

	return p
}

func (p Pool) run(f func(string)) {
	for available := p.size; ; {
		select {
		case <-p.removeJob:
			available++
		case path := <-p.addJob:
			p.jobs = append(p.jobs, path)
		}

		if available > 0 && len(p.jobs) > 0 {
			go f(p.jobs[0])
			p.jobs = p.jobs[1:]
			available--
			continue
		}

		// make sure there are no pending jobs before closing, since the order of select isn't deterministic
		if available == p.size {
			select {
			case path := <-p.addJob:
				p.jobs = append(p.jobs, path)
				continue
			default:
			}

			break
		}
	}

	p.wait <- nil
}
