module example.com/main

go 1.16

replace example.com/graph => ../graph

replace example.com/handler => ../handler

require (
	example.com/graph v0.0.0-00010101000000-000000000000
	example.com/handler v0.0.0-00010101000000-000000000000
)
