package main

type Method func(c *callImpl)

func GET() Method {
	return func(c *callImpl) {
		c.gomockCall = c.gomockCall.Do(func(path string) {
			c.commitCall()
		})
	}
}
