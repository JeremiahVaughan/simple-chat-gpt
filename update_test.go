package main

import (
	"os"
	"strings"
	"testing"
)

func Test_applyWordWrap(t *testing.T) {
	t.Run("easy mode", func(t *testing.T) {
		words := `happy duck paste`
		expected := `happy duck paste`
		width := 40
		got := applyWordWrap(words, width)
		if got != expected {
			t.Errorf("error, expected: %sEND, but got: %sEND", strings.ReplaceAll(expected, " ", "-"), strings.ReplaceAll(got, " ", "-"))
		}
	})

	t.Run("hard mode", func(t *testing.T) {
		words := `happy duck paste mustache pottery potluckz`
		expected := `happy duck paste
mustache pottery
potluckz`
		width := 20
		got := applyWordWrap(words, width)
		if got != expected {
			t.Errorf("error, expected: %s, but got: %s", strings.ReplaceAll(expected, " ", "-"), strings.ReplaceAll(got, " ", "-"))
		}
	})

	t.Run("harder mode", func(t *testing.T) {
		words := `package main

import "fmt"
	
func main() {
    fmt.Println("Hello, World!")
}`
		expected := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`
		width := 80
		got := applyWordWrap(words, width)
		if got != expected {
			t.Errorf("error, expected: %s, but got: %s", expected, got)
		}
	})

	t.Run("hardest mode", func(t *testing.T) {
		words := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World there cruel world, I want to let you know that I don't like it!")
}`
		expected := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World there cruel world, I want to let you know that I
don't like it!")
}`
		width := 80
		got := applyWordWrap(words, width)
		if got != expected {
			t.Errorf("error, expected: %sEND, but got: %sEND", strings.ReplaceAll(expected, " ", "-"), strings.ReplaceAll(got, " ", "-"))
		}
	})

	t.Run("indented list", func(t *testing.T) {
		words := `
List:
    1. Eggs
    2. Bacon`
		expected := `
List:
    1. Eggs
    2. Bacon`
		width := 80
		got := applyWordWrap(words, width)
		if got != expected {
			t.Errorf("error, expected: %sEND, but got: %sEND", strings.ReplaceAll(expected, " ", "-"), strings.ReplaceAll(got, " ", "-"))
		}
	})

	t.Run("edge case should never exceed width", func(t *testing.T) {
		words := `pack main
some fluff
ok`
		expected := `pack main
some
fluff
ok`
		width := 10
		got := applyWordWrap(words, width)
		if got != expected {
			t.Errorf("error, expected: %s, but got: %s", strings.ReplaceAll(expected, " ", "-"), strings.ReplaceAll(got, " ", "-"))
		}
	})

	t.Run("edge case should never exceed width", func(t *testing.T) {
		words := `
package main                                                          

import (                                                              
    "fmt"                                                                 
)                                                                     

func sillyDuckCountdown() {                                           
    snacks := []string{"quackers", "bread crumbs", "fish flakes", "corn", "ice cream"}                                                          
    
    fmt.Println("🐥 Quack! Let's count down my favorite snacks!")         
    
    for i := len(snacks); i > 0; i-- {                                    
        fmt.Printf("%d... %s!\n", i, snacks[i-1])                             
    }                                                                     
    
    fmt.Println("🎉 Yay! All my snacks are here! Time to eat! 🍽️")        
}                                                                     

func main() {                                                         
    sillyDuckCountdown()                                                  
}                                                                     
`
		expected := `
package main

import (
    "fmt"
)

func sillyDuckCountdown() {
    snacks := []string{"quackers", "bread crumbs", "fish flakes", "corn",
"ice cream"}
    
    fmt.Println("🐥 Quack! Let's count down my favorite snacks!")
    
    for i := len(snacks); i > 0; i-- {
        fmt.Printf("%d... %s!\n", i, snacks[i-1])
    }
    
    fmt.Println("🎉 Yay! All my snacks are here! Time to eat! 🍽️")
}

func main() {
    sillyDuckCountdown()
}
`
		width := 75
		got := applyWordWrap(words, width)
		if got != expected {
			t.Errorf("error, expected: %s, but got: %s", strings.ReplaceAll(expected, " ", "-"), strings.ReplaceAll(got, " ", "-"))
		}
	})

	t.Run("ultimate test", func(t *testing.T) {
		input := `**Title: The Echo of Stars** 

In a quaint village nestled between lush green hills and shimmering lakes, there lived a young girl named Elara. She was known for her insatiable curiosity and playful spirit. Every evening, she would sit outside her little cottage, gazing up at the radiant tapestry of stars that blanketed the night sky. Her grandmother, a wise woman  with silver hair and a twinkle in her eye, often shared tales of the constellations and the ancient legends that surrounded them.

One evening, as twilight ushered in the first hints of starlight, Elara’s grandmother whispered of a special star—the Celestia Solara. Legends spoke of its power to grant one wish to those with a pure heart. The moment the star blinked its light three times, a wish would come true. However, it was said that the star was hidden deep within the Enchanted Forest, guarded by the Spirit of the Night.      

That night, feeling a sense of adventure surge within her, Elara vowed to find this mystical star. She gathered her essentials—a small lantern, a map drawn by her grandmother, and her favorite scarf—and set out into the waiting arms of the forest, guided by the soft glow of fireflies.

The forest was alive with sounds, with rustling leaves and distant hoots guiding her deeper into its embrace. As she walked, she encountered various magical creatures—a wise old owl who spoke in riddles, a mischievous fox who tried to lead her astray, and a gentle deer that offered her directions. Each encounter taught her more about the world and herself, filling her heart with courage and wonder.

After hours of wandering, Elara reached a clearing where the trees parted to reveal a silver pond, its surface reflecting the starlit sky above. In the center of the pond floated a single shimmering lily pad, upon which sat the Spirit of the Night, draped in a cloak of midnight clouds sprinkled with stardust. Its voice echoed like soft chimes as it spoke, "To find Celestia Solara, you must answer a riddle: What is the greatest treasure, yet cannot be held?"

Elara thought deeply. She considered wealth, love, and knowledge. Finally, she realized the answer was simple yet profound: "It is hope."

With a nod, the Spirit smiled and gestured towards the sky. In a moment, the stars began to swirl, and there it was—Celestia Solara! It twinkled brightly, pulsing with anticipation. Elara closed her eyes and made her wish, heartfelt and sincere. She wished for her village to always be filled with light and laughter, for joy to resonate in every heart.

As she opened her eyes, the star blinked thrice, showering the forest with a cascade of shimmering light. She felt a warmth envelop her, and when the glow faded, she found herself back outside her cottage, the morning sun rising in the sky.

The villagers awakened to a new day, filled with laughter and tales of wonder. The echoes of Elara’s wish lingered, and from that day forth, joy became a part of their everyday lives.

Every night, when Elara gazed at the stars, she remembered the adventure and the magic of her wish. She learned that while she could not hold hope in her hands, she could carry it in her heart, and share it with everyone around her.

And so, the village flourished, forever illuminated by the echoes of stars and the warmth of a girl’s pure wish.`
		expected := `**Title: The Echo of Stars**

In a quaint village nestled between lush green hills and shimmering lakes,
there lived a young girl named Elara. She was known for her insatiable
curiosity and playful spirit. Every evening, she would sit outside her
little cottage, gazing up at the radiant tapestry of stars that blanketed
the night sky. Her grandmother, a wise woman with silver hair and a
twinkle in her eye, often shared tales of the constellations and the
ancient legends that surrounded them.

One evening, as twilight ushered in the first hints of starlight,
Elara’s grandmother whispered of a special star—the Celestia Solara.
Legends spoke of its power to grant one wish to those with a pure heart.
The moment the star blinked its light three times, a wish would come true.
However, it was said that the star was hidden deep within the Enchanted
Forest, guarded by the Spirit of the Night.

That night, feeling a sense of adventure surge within her, Elara vowed to
find this mystical star. She gathered her essentials—a small lantern, a
map drawn by her grandmother, and her favorite scarf—and set out into
the waiting arms of the forest, guided by the soft glow of fireflies.

The forest was alive with sounds, with rustling leaves and distant hoots
guiding her deeper into its embrace. As she walked, she encountered
various magical creatures—a wise old owl who spoke in riddles, a
mischievous fox who tried to lead her astray, and a gentle deer that
offered her directions. Each encounter taught her more about the world and
herself, filling her heart with courage and wonder.

After hours of wandering, Elara reached a clearing where the trees parted
to reveal a silver pond, its surface reflecting the starlit sky above. In
the center of the pond floated a single shimmering lily pad, upon which
sat the Spirit of the Night, draped in a cloak of midnight clouds
sprinkled with stardust. Its voice echoed like soft chimes as it spoke,
"To find Celestia Solara, you must answer a riddle: What is the greatest
treasure, yet cannot be held?"

Elara thought deeply. She considered wealth, love, and knowledge. Finally,
she realized the answer was simple yet profound: "It is hope."

With a nod, the Spirit smiled and gestured towards the sky. In a moment,
the stars began to swirl, and there it was—Celestia Solara! It twinkled
brightly, pulsing with anticipation. Elara closed her eyes and made her
wish, heartfelt and sincere. She wished for her village to always be
filled with light and laughter, for joy to resonate in every heart.

As she opened her eyes, the star blinked thrice, showering the forest with
a cascade of shimmering light. She felt a warmth envelop her, and when the
glow faded, she found herself back outside her cottage, the morning sun
rising in the sky.

The villagers awakened to a new day, filled with laughter and tales of
wonder. The echoes of Elara’s wish lingered, and from that day forth,
joy became a part of their everyday lives.

Every night, when Elara gazed at the stars, she remembered the adventure
and the magic of her wish. She learned that while she could not hold hope
in her hands, she could carry it in her heart, and share it with everyone
around her.

And so, the village flourished, forever illuminated by the echoes of stars
and the warmth of a girl’s pure wish.`
		width := 75
		got := applyWordWrap(input, width)
		os.WriteFile("test-ruselts", []byte(got), 0777)
		if expected != got {
			t.Errorf("error, expected: %s, but got: %s", expected, got)
		}
	})
}
