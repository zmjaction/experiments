module example.com/lazy

go 1.17

require example.com/a v0.1.0

require (
	example.com/b v0.1.0 // indirect
	example.com/c v0.1.0 // indirect
)

replace (
	example.com/a v0.1.0 => ./a
	example.com/b v0.1.0 => ./b
	example.com/c v0.1.0 => ./c1
	example.com/c v0.2.0 => ./c2
)
