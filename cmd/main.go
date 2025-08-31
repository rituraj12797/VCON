package main

import (
	"fmt"
	"strings" // Import the strings package
	"vcon/internal/api"
	"vcon/internal/engine"
	"vcon/internal/globalStore"
	"vcon/internal/storage"
)

func main() {
	fmt.Println(" hello world ")

	api.DemoAPI()
	api.DemoHandler()
	x := storage.NewTree()

	// for now we define global variable in their files thwmselvs only

	globalStore.Initialize()

	gs := globalStore.GlobalStore

	// now try converting to " A mango can run but this can not "

	x.AddNode(0, 1, "base version", nil)
	x.AddNode(1, 2, "new version", nil)
	x.AddNode(1, 3, "just new", nil)
	x.AddNode(3, 4, "gotu", nil)

	arr := []int{4, 2, 4, 7, 9, 5, 1, 0, 4, 2}
	brr := []int{4, 8, 9, 4, 10, 7, 6, 5, 0, 1}
	crr := []int{4, 2, 7, 8, 6, 4, 1, 2, 4, 5, 6}

	y, _ := engine.LCS(&arr, &brr)
	k, _ := engine.LCS(&brr, &crr)

	z := engine.GenerateDelta(&arr, &brr, &y)
	d := engine.GenerateDelta(&brr, &crr, &k)

	c := engine.ApplyDelta(arr, z)
	f := engine.ApplyDelta(c, d)

	// x.ShowTree()
	fmt.Println(" The delta is : ", z)
	fmt.Println(" The CRR  :       ", crr)
	fmt.Println(" The resultant  : ", f)

	// --- Document Test for ContentRenderer ---
	fmt.Println("\n--- Running Document Test for ContentRenderer ---")

	story := `Title: Bimbo the Button-Nosed Bear
Once upon a humming morning in Bluebell Forest, the sun stretched golden fingers across a patchwork of flowers.
Among the fluttering butterflies lived Bimbo, a plump little cartoon bear with a nose as round and shiny as a black button and ears that wiggled when he laughed.
Chapter 1: The Lost Melody
Bimbo loved collecting sounds.
Every day he roamed the forest with an acorn-shell satchel, scooping up giggles, bird chirps, and rustling leaves to keep in glimmering jars at home.
One morning he noticed the forest was suspiciously quiet.
The robins refused to sing, the crickets had zipped their legs shut, and even the brook flowed in gentle silence.
Bimbo’s satchel lay empty—something precious was missing: the Heart Melody, the tune that stitched every other sound together.
Chapter 2: A Map of Echoes
Determined, Bimbo visited his best friend, Lulu the wise ladybug librarian.
Lulu perched on a toadstool and presented a tiny parchment: the Map of Echoes.
“Follow the dotted hushes,” she said, pointing to faint silver lines that appeared only when you listened carefully.
Bimbo tucked the map into his satchel and set off, his button nose twitching with excitement.
Chapter 3: Trials of Quiet
First came Whispering Meadow, where drifting dandelion seeds tried to distract him with sleepy sighs.
Bimbo hummed softly, refusing to forget the missing melody.
Next, he crossed Murmur Creek on polished pebble stepping-stones.
Each stone silenced the splash beneath it, but Bimbo tapped them in a steady rhythm to keep hope alive.
Chapter 4: The Silent Peak
At dusk he reached Silent Peak, a hill so still the clouds tip-toed overhead.
There he found a lonely Windharp—an ancient stone arch strung with silver vines—and behind it stood a shy Moon Moth, its wings folded tight.
The moth had accidentally tangled its shimmering antennae in the vines, muffling every breeze.
Without wind, no harp could sing, and thus the Heart Melody had vanished.
Chapter 5: A Gentle Rescue
Bimbo spoke softly.
“May I help?”
The Moon Moth blinked moonlit eyes and nodded.
With careful claws, Bimbo freed the silk-fine antennae.
As soon as the last thread loosened, a breeze tumbled through the harp.
Strings glowed, and the Heart Melody—bright, warm, and full—washed over the forest like sunrise after rain.
Chapter 6: A Symphony Returns
On their journey back, each place they passed blossomed with sound.
River splashes harmonized with chirping robins; dandelion seeds danced to cricket fiddles; even Lulu’s library chimed with rustling pages.
Bimbo filled jar after jar with joyous notes, but saved the Heart Melody for everyone to share.
Epilogue: The Ever-Full Satchel
That night, under firefly lanterns, Bluebell Forest held the first Festival of Echoes.
Creatures tiny and tall sang together, and Bimbo learned a wonderful secret: the more you share a melody, the louder your own heart sings.
From then on, Bimbo’s satchel never felt empty again, for kindness—and music—have a way of coming home to any bear with a button nose and listening ears.
The End`

	// Split the story into lines
	lines := strings.Split(story, "\n")

	// Intern each line and collect the IDs
	var storyIDs []int
	for _, line := range lines {
		if line != "" { // Avoid interning empty lines
			id, err := gs.Intern(line)
			if err != nil {
				fmt.Printf("Error interning line: %v\n", err)
				continue
			}
			storyIDs = append(storyIDs, id)
		}
	}

	// Pass the document through the content renderer
	engine.ContentRendered(&storyIDs)

}
