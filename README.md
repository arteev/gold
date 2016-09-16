# gold

[![Build Status](https://travis-ci.org/arteev/gold.svg?branch=master)](https://travis-ci.org/arteev/gold)


Library Golang Line Display


Description
-----------

Library line displays based on Firich,Datecs (serial port)


Supported 
---------    
    [+] Firich Command (partially)
    [ ] Datecs Command


Quick start
-----------
```
	config := map[string]interface{}{
		"Name": "/dev/ttyUSB0",
		"Baud": 9600,
	}
	linedsp, err := display.Open("serial", config)
	if err != nil {
		log.Fatal(err)
	}
	dsp, err := linedsp.GetDisplay(&firich.FirichProtocol{})
	if err != nil {
		log.Fatal(err)
	}
	defer dsp.Close()
	dsp.SetEncoding(charmap.CodePage866)
	if err := dsp.Init(); err != nil {
		fmt.Println(err)
	}
	dsp.PrintRow(1, "Price:10$ Quant:2")
	dsp.PrintRow(2, "Total:20$")

```

License
-------

  MIT

Author
------

Arteev Aleksey
