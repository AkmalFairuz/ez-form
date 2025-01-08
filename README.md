# ez-form

Form library for dragonfly server

## Example

```go
func exampleMenu(p *player.Player) {
	m := form.NewMenu().WithTitle("Server Selector")
	m.WithContent("Available servers:")
	m.WithButton("Lobby-1", "textures/items/compass")
	m.WithButton("Lobby-2", "textures/items/compass")
	m.WithButton("Lobby-3", "textures/items/compass")
	m.WithCallback(func(p *player.Player, result int) {
		fmt.Printf("Player %s selected server Lobby-%d\n", p.Name(), result+1)
	})
	m.WithCloseCallback(func(p *player.Player) {
		fmt.Println("Menu closed")
	})
	p.SendForm(m)
}

func exampleModal(p *player.Player) {
	m := form.NewModal().WithTitle("Confirmation")
	m.WithContent("Are you sure you want to delete this item?")
	m.WithButton1("Yes")
	m.WithButton2("No")
	m.WithCallback(func(p *player.Player, button1 bool) {
		if button1 {
			fmt.Println("Player confirmed deletion")
		} else {
			fmt.Println("Player cancelled deletion")
		}
	})
	m.WithCloseCallback(func(p *player.Player) {
		fmt.Println("modal closed")
	})
	p.SendForm(m)
}

func exampleCustom(p *player.Player) {
	c := form.NewCustom().WithTitle("Custom Form")
	c.WithElement("name", form.NewInput("Your name").WithPlaceholder("Pig"))
	c.WithElement("age", form.NewSlider("Your age", 0, 200).WithStepSize(1).WithDefault(10))
	c.WithElement("favourite color", form.NewDropdown("Select your favorite color").WithOptions("Red", "Green", "Blue"))
	c.WithElement("something", form.NewToggle("Enable something", false))
	c.WithCallback(func(p *player.Player, response form.CustomResponse) {
		var data struct {
			Name         string `form:"name"`
			Age          int    `form:"age"`
			FavouriteCol string `form:"favourite color"`
			Something    bool   `form:"something"`
		}

		if err := response.Bind(&data); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Player %s submitted form with name %s, age %d, color %s and something %t\n", p.Name(), data.Name, data.Age, data.FavouriteCol, data.Something)

		// Manual parsing
		// name := response.String("name")
		// age := response.Int("age")
		// colorIndex := response.Int("favourite color")
		// something := response.Bool("something")
		// fmt.Printf("Player %s submitted form with name %s, age %d, color %s and something %t\n", p.Name(), name, age, colorIndex, something)
	})
	c.WithCloseCallback(func(p *player.Player) {
		fmt.Println("Custom form closed")
	})
	p.SendForm(c)
}

```