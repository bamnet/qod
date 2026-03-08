package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type quote struct {
	text   string
	author string
}

// quotes is the pool of daily quotes, rotated by day-of-year.
// Add more entries freely — order doesn't matter.
var quotes = []quote{
	// Philosophers
	{text: "The unexamined life is not worth living.", author: "Socrates"},
	{text: "The only true wisdom is in knowing you know nothing.", author: "Socrates"},
	{text: "We are what we repeatedly do. Excellence, then, is not an act, but a habit.", author: "Aristotle"},
	{text: "Happiness depends upon ourselves.", author: "Aristotle"},
	{text: "The energy of the mind is the essence of life.", author: "Aristotle"},
	{text: "Courage is knowing what not to fear.", author: "Plato"},
	{text: "A room without books is like a body without a soul.", author: "Marcus Tullius Cicero"},
	{text: "Not how long, but how well you have lived is the main thing.", author: "Seneca"},
	{text: "Give me a lever long enough and a fulcrum on which to place it, and I shall move the world.", author: "Archimedes"},
	{text: "I think, therefore I am.", author: "René Descartes"},
	{text: "That which does not kill us makes us stronger.", author: "Friedrich Nietzsche"},
	{text: "He who has a why to live can bear almost any how.", author: "Friedrich Nietzsche"},
	{text: "The good life is one inspired by love and guided by knowledge.", author: "Bertrand Russell"},
	{text: "Man is condemned to be free.", author: "Jean-Paul Sartre"},
	{text: "In the depth of winter, I finally learned that within me there lay an invincible summer.", author: "Albert Camus"},
	{text: "What is to give light must endure burning.", author: "Viktor Frankl"},
	{text: "Everything can be taken from a man but one thing: the last of human freedoms — to choose one's attitude.", author: "Viktor Frankl"},

	// Scientists
	{text: "Imagination is more important than knowledge.", author: "Albert Einstein"},
	{text: "Two things are infinite: the universe and human stupidity; and I'm not sure about the universe.", author: "Albert Einstein"},
	{text: "In the middle of difficulty lies opportunity.", author: "Albert Einstein"},
	{text: "I have no special talent. I am only passionately curious.", author: "Albert Einstein"},
	{text: "Logic will get you from A to B. Imagination will take you everywhere.", author: "Albert Einstein"},
	{text: "Nothing in life is to be feared, it is only to be understood.", author: "Marie Curie"},
	{text: "If I have seen further, it is by standing on the shoulders of giants.", author: "Isaac Newton"},
	{text: "All truths are easy to understand once they are discovered; the point is to discover them.", author: "Galileo Galilei"},
	{text: "And yet it moves.", author: "Galileo Galilei"},
	{text: "Genius is one percent inspiration, ninety-nine percent perspiration.", author: "Thomas Edison"},
	{text: "I have not failed. I've just found 10,000 ways that won't work.", author: "Thomas Edison"},
	{text: "The present is theirs; the future, for which I really worked, is mine.", author: "Nikola Tesla"},
	{text: "If you wish to make an apple pie from scratch, you must first invent the universe.", author: "Carl Sagan"},
	{text: "Somewhere, something incredible is waiting to be known.", author: "Carl Sagan"},

	// Eastern thinkers
	{text: "It does not matter how slowly you go as long as you do not stop.", author: "Confucius"},
	{text: "Wherever you go, go with all your heart.", author: "Confucius"},
	{text: "Our greatest glory is not in never falling, but in rising every time we fall.", author: "Confucius"},
	{text: "The supreme art of war is to subdue the enemy without fighting.", author: "Sun Tzu"},
	{text: "Know your enemy and know yourself, and you will never be defeated.", author: "Sun Tzu"},
	{text: "The wound is the place where the Light enters you.", author: "Rumi"},
	{text: "Out beyond ideas of wrongdoing and rightdoing, there is a field. I'll meet you there.", author: "Rumi"},
	{text: "You can't cross the sea merely by standing and staring at the water.", author: "Rabindranath Tagore"},

	// Writers & Poets
	{text: "Some are born great, some achieve greatness, and some have greatness thrust upon them.", author: "William Shakespeare"},
	{text: "All the world's a stage, and all the men and women merely players.", author: "William Shakespeare"},
	{text: "To thine own self be true.", author: "William Shakespeare"},
	{text: "We know what we are, but know not what we may be.", author: "William Shakespeare"},
	{text: "Beauty is truth, truth beauty.", author: "John Keats"},
	{text: "It is better to have loved and lost than never to have loved at all.", author: "Alfred Lord Tennyson"},
	{text: "If you tell the truth, you don't have to remember anything.", author: "Mark Twain"},
	{text: "The secret of getting ahead is getting started.", author: "Mark Twain"},
	{text: "All you need in this life is ignorance and confidence, and then success is sure.", author: "Mark Twain"},
	{text: "The two most important days in your life are the day you are born and the day you find out why.", author: "Mark Twain"},
	{text: "Hope is the thing with feathers that perches in the soul.", author: "Emily Dickinson"},
	{text: "Not all those who wander are lost.", author: "J.R.R. Tolkien"},
	{text: "I took the road less traveled by, and that has made all the difference.", author: "Robert Frost"},
	{text: "I am large, I contain multitudes.", author: "Walt Whitman"},
	{text: "All happy families are alike; each unhappy family is unhappy in its own way.", author: "Leo Tolstoy"},
	{text: "To live is the rarest thing in the world. Most people exist, that is all.", author: "Oscar Wilde"},
	{text: "Be yourself; everyone else is already taken.", author: "Oscar Wilde"},
	{text: "We are all in the gutter, but some of us are looking at the stars.", author: "Oscar Wilde"},
	{text: "I am not afraid of storms, for I am learning how to sail my ship.", author: "Louisa May Alcott"},
	{text: "Tell me, what is it you plan to do with your one wild and precious life?", author: "Mary Oliver"},

	// Political & civic leaders
	{text: "The only thing we have to fear is fear itself.", author: "Franklin D. Roosevelt"},
	{text: "Great minds discuss ideas; average minds discuss events; small minds discuss people.", author: "Eleanor Roosevelt"},
	{text: "The future belongs to those who believe in the beauty of their dreams.", author: "Eleanor Roosevelt"},
	{text: "In any moment of decision, the best thing you can do is the right thing.", author: "Theodore Roosevelt"},
	{text: "Do what you can, with what you have, where you are.", author: "Theodore Roosevelt"},
	{text: "Whatever you are, be a good one.", author: "Abraham Lincoln"},
	{text: "Government of the people, by the people, for the people, shall not perish from the earth.", author: "Abraham Lincoln"},
	{text: "In matters of style, swim with the current; in matters of principle, stand like a rock.", author: "Thomas Jefferson"},
	{text: "Give me liberty, or give me death!", author: "Patrick Henry"},
	{text: "Ask not what your country can do for you — ask what you can do for your country.", author: "John F. Kennedy"},
	{text: "The price of greatness is responsibility.", author: "Winston Churchill"},
	{text: "Success is not final, failure is not fatal: it is the courage to continue that counts.", author: "Winston Churchill"},
	{text: "Now this is not the end. It is, perhaps, the end of the beginning.", author: "Winston Churchill"},
	{text: "Never interrupt your enemy when he is making a mistake.", author: "Napoleon Bonaparte"},
	{text: "Impossible is a word found only in the dictionary of fools.", author: "Napoleon Bonaparte"},
	{text: "I learned that courage was not the absence of fear, but the triumph over it.", author: "Nelson Mandela"},
	{text: "It always seems impossible until it's done.", author: "Nelson Mandela"},

	// Social leaders
	{text: "Darkness cannot drive out darkness; only light can do that. Hate cannot drive out hate; only love can do that.", author: "Martin Luther King Jr."},
	{text: "The time is always right to do what is right.", author: "Martin Luther King Jr."},
	{text: "Injustice anywhere is a threat to justice everywhere.", author: "Martin Luther King Jr."},
	{text: "An eye for an eye will only make the whole world blind.", author: "Mahatma Gandhi"},
	{text: "Be the change that you wish to see in the world.", author: "Mahatma Gandhi"},
	{text: "Live as if you were to die tomorrow. Learn as if you were to live forever.", author: "Mahatma Gandhi"},
	{text: "To be yourself in a world that is constantly trying to make you something else is the greatest accomplishment.", author: "Ralph Waldo Emerson"},
	{text: "Do not go where the path may lead, go instead where there is no path and leave a trail.", author: "Ralph Waldo Emerson"},
	{text: "Life is either a daring adventure or nothing at all.", author: "Helen Keller"},
	{text: "The most courageous act is still to think for yourself. Aloud.", author: "Coco Chanel"},
	{text: "Float like a butterfly, sting like a bee.", author: "Muhammad Ali"},
	{text: "The man who has no imagination has no wings.", author: "Muhammad Ali"},
}

type content struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Text     string `json:"text"`
	RenderAs string `json:"render_as"`
	Duration int    `json:"duration,omitempty"`
}

func todaysQuote() quote {
	now := time.Now()
	// Use a fixed seed based on year+day so the quote is stable all day
	// but rotates through the full pool across days.
	day := now.Year()*1000 + now.YearDay()
	return quotes[day%len(quotes)]
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	q := todaysQuote()
	text := fmt.Sprintf("\u201c%s\u201d \u2014 %s", q.text, q.author)

	feed := []content{
		{
			Type:     "RichText",
			Name:     "Quote of the Day",
			Text:     text,
			RenderAs: "plaintext",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(feed); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/feed", feedHandler)

	log.Printf("Serving quote of the day on :%s/feed", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
