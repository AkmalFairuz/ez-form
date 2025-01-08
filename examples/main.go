package main

import (
	"fmt"
	form "github.com/akmalfairuz/ez-form"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/pelletier/go-toml"
	"log/slog"
	"os"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	chat.Global.Subscribe(chat.StdoutSubscriber{})
	conf, err := readConfig(slog.Default())
	if err != nil {
		panic(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()

	srv.Listen()
	for p := range srv.Accept() {
		_ = p
		p.Handle(exampleHandler{})
		_ = p.Inventory().SetItem(0, item.NewStack(item.DragonBreath{}, 1).WithCustomName("Custom Form").WithValue("form", "custom"))
		_ = p.Inventory().SetItem(1, item.NewStack(item.DragonBreath{}, 1).WithCustomName("Modal Form").WithValue("form", "modal"))
		_ = p.Inventory().SetItem(2, item.NewStack(item.DragonBreath{}, 1).WithCustomName("Menu Form").WithValue("form", "menu"))
	}
}

type exampleHandler struct {
	player.NopHandler
}

func (h exampleHandler) HandleItemUse(ctx *player.Context) {
	it, _ := ctx.Val().HeldItems()
	itemVal, _ := it.Value("form")
	switch itemVal {
	case "custom":
		exampleCustom(ctx.Val())
	case "modal":
		exampleModal(ctx.Val())
	case "menu":
		exampleMenu(ctx.Val())
	}
}

func exampleMenu(p *player.Player) {
	m := form.NewMenu("Server Selector")
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
	m := form.NewModal("Confirmation")
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
	c := form.NewCustom("Custom Form")
	c.WithElement("name", form.NewInput("Your name").WithPlaceholder("Pig"))
	c.WithElement("age", form.NewSlider("Your age", 0, 200).WithStepSize(1).WithDefault(10))
	c.WithElement("favourite color", form.NewDropdown("Select your favorite color").WithOptions("Red", "Green", "Blue"))
	c.WithElement("something", form.NewToggle("Enable something", false))
	c.WithElement("num", form.NewStepSlider("Select a number").WithOptions("1", "2", "3", "4", "5"))
	c.WithCallback(func(p *player.Player, response form.CustomResponse) {
		var data struct {
			Name        string `form:"name"`
			Age         int    `form:"age"`
			FavColorIdx int    `form:"favourite color"`
			Something   bool   `form:"something"`
			Num         int    `form:"num"`
		}

		if err := response.Bind(&data); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Player %s submitted form with name %s, age %d, colorIndex %d, something %t, and num %d\n", p.Name(), data.Name, data.Age, data.FavColorIdx, data.Something, data.Num+1)

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

// readConfig reads the configuration from the config.toml file, or creates the
// file if it does not yet exist.
func readConfig(log *slog.Logger) (server.Config, error) {
	c := server.DefaultConfig()
	var zero server.Config
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return zero, fmt.Errorf("encode default config: %v", err)
		}
		if err := os.WriteFile("config.toml", data, 0644); err != nil {
			return zero, fmt.Errorf("create default config: %v", err)
		}
		return c.Config(log)
	}
	data, err := os.ReadFile("config.toml")
	if err != nil {
		return zero, fmt.Errorf("read config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return zero, fmt.Errorf("decode config: %v", err)
	}
	return c.Config(log)
}
