package runner

import "fmt"

var banner = fmt.Sprintf(`       ______________                   /\           /\   /—/
      |______________|                  \_\          \_\  \_/
           ||  ||                    /\________   _______________
    _______||__||_______             \______/  /  \_____________/
   |_______||__||_______|                 /  /          / /
          / /  \ \                    __/__/____   ||=========||
        /  /    \  \                 |______/  /   ||=========||	
      /   /      \   \                    /  /     ||=========||
    /    /        \    \                /  /       ||=========||
  /     /          \     \            /\ _/\_________________________/\
/______/            \______\          \__\__________________________/_/
`)

func showBanner() {
	fmt.Printf("%s\n", banner)
}
