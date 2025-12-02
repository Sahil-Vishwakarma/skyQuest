package hints

import (
	"math/rand"
	"sync"
	"time"
)

var (
	// Cache for city facts to avoid repeated lookups
	factCache   = make(map[string]string)
	factCacheMu sync.RWMutex

	// City-specific facts mapped by city name (facts don't mention the city name)
	cityFacts = map[string][]string{
		// USA Cities
		"New York": {
			"This city has over 800 languages spoken, making it the most linguistically diverse place on Earth.",
			"The subway system here has 472 stations, the most of any metro system in the world.",
			"Central Park receives around 42 million visitors annually, more than most countries.",
			"The iconic yellow taxis here number over 13,000 vehicles.",
		},
		"Newark": {
			"This city is home to the oldest museum in the state, founded in 1909.",
			"The local airport was the first major airport in the metropolitan area, opening in 1928.",
			"Branch Brook Park here has more cherry blossom trees than Washington D.C.",
		},
		"Los Angeles": {
			"The Hollywood sign was originally built in 1923 to advertise a real estate development.",
			"This city produces more than 10 million pounds of avocados annually.",
			"The entertainment industry here generates over $50 billion in annual revenue.",
			"This is the only city to host the Summer Olympics twice in the USA.",
		},
		"San Francisco": {
			"The famous cable cars here are the only mobile National Historic Landmark in the USA.",
			"The Golden Gate Bridge contains enough wire to circle the equator three times.",
			"This city is built on more than 50 hills.",
			"Fortune cookies were actually invented here, not in China.",
		},
		"Chicago": {
			"The first-ever skyscraper was built here in 1885.",
			"The river here flows backwards - engineers reversed it in 1900.",
			"Deep-dish pizza was invented in this city in 1943.",
			"This city has the largest collection of Impressionist paintings outside of Paris.",
		},
		"Miami": {
			"This is the only major U.S. city founded by a woman.",
			"The Art Deco Historic District here has over 800 preserved buildings.",
			"This city has more than 12 miles of beaches.",
			"Spanish is the first language for over 60% of residents here.",
		},
		"Boston": {
			"America's first public park, first public school, and first subway were all built here.",
			"The local accent drops the letter 'R' - a linguistic quirk from early British settlers.",
			"Harvard University, founded in 1636, is the oldest university in the USA.",
			"The Boston Marathon is the world's oldest annual marathon.",
		},
		"Atlanta": {
			"This city is home to the world's busiest airport by passenger traffic.",
			"Coca-Cola was invented here in 1886.",
			"This is the only major American city destroyed during the Civil War.",
			"The Martin Luther King Jr. National Historic Site is located here.",
		},
		"Dallas": {
			"The frozen margarita machine was invented here in 1971.",
			"This city has more restaurants per capita than any other U.S. city.",
			"The State Fair here features a 55-foot-tall cowboy statue named Big Tex.",
			"JFK was assassinated in this city in 1963.",
		},
		"Seattle": {
			"The first Starbucks opened here in 1971 at Pike Place Market.",
			"This city gets less annual rainfall than New York, Miami, or Houston.",
			"The Space Needle was built for the 1962 World's Fair.",
			"This city is home to the world's first gas station, opened in 1907.",
		},
		// Canada
		"Toronto": {
			"The CN Tower was the world's tallest free-standing structure for 34 years.",
			"This city has the largest underground shopping complex in the world - PATH.",
			"Over 140 languages are spoken here, making it one of the most diverse cities.",
			"The local film industry is the third largest in North America after LA and NYC.",
		},
		"Vancouver": {
			"This city has been ranked the most livable city in the world multiple times.",
			"Stanley Park here is larger than New York's Central Park.",
			"No building in the downtown core can be taller than the nearby mountains.",
			"The 'Hollywood North' nickname comes from filming over 65 movies annually.",
		},
		// Mexico
		"Mexico City": {
			"This city is sinking at a rate of 10 inches per year due to over-extraction of groundwater.",
			"It was built on the ruins of the ancient Aztec capital Tenochtitlan.",
			"The metro system here has unique icons for each station to help illiterate riders.",
			"This is the oldest capital city in the Americas.",
		},
		// UK
		"London": {
			"Big Ben is actually the name of the bell, not the clock tower.",
			"The Underground here is the oldest metro system in the world, opened in 1863.",
			"There are over 170 museums here, many of which are free.",
			"The city has been the capital of seven different kingdoms throughout history.",
		},
		// France
		"Paris": {
			"The Eiffel Tower was supposed to be dismantled after 20 years.",
			"There's only one stop sign in the entire city.",
			"The Louvre would take 100 days to see everything if you spent 30 seconds on each piece.",
			"This city has 450 parks and gardens.",
		},
		// Germany
		"Frankfurt": {
			"This city is home to the European Central Bank.",
			"The local airport is the largest in Germany and a major European hub.",
			"Goethe, Germany's most famous writer, was born here.",
			"The skyline is nicknamed 'Mainhattan' due to its high-rise buildings.",
		},
		"Munich": {
			"Oktoberfest started here in 1810 as a royal wedding celebration.",
			"The local BMW headquarters is shaped like a four-cylinder engine.",
			"This city has one of the largest urban parks in the world - the English Garden.",
			"The famous Glockenspiel clock has been entertaining crowds since 1908.",
		},
		// Netherlands
		"Amsterdam": {
			"This city has more bicycles than people - about 881,000 bikes.",
			"There are 165 canals totaling over 100 kilometers.",
			"The houses are narrow because property tax was once based on building width.",
			"The Rijksmuseum houses over 1 million artworks.",
		},
		// Spain
		"Madrid": {
			"This is the highest capital city in Europe at 667 meters above sea level.",
			"The Prado Museum has one of the world's finest collections of European art.",
			"Dinner here typically starts at 10 PM, among the latest in Europe.",
			"The city's symbol is a bear and a strawberry tree.",
		},
		"Barcelona": {
			"The Sagrada Familia has been under construction for over 140 years.",
			"This city has 4.5 km of beaches within the city limits.",
			"La Rambla street is one of the most famous pedestrian streets in the world.",
			"The architect Gaudí designed many of the city's most famous buildings.",
		},
		// Italy
		"Rome": {
			"There's a country entirely within this city - Vatican City.",
			"Visitors throw about €3,000 into the Trevi Fountain daily.",
			"This city has more ancient obelisks than any other place in the world.",
			"The Pantheon has been in continuous use for 2,000 years.",
		},
		// Switzerland
		"Zurich": {
			"The Swiss banking industry manages about $7 trillion in assets here.",
			"This city was ranked the world's best city for quality of life multiple times.",
			"Einstein developed his theory of special relativity while living here.",
			"The local lake provides drinking water directly from the tap.",
		},
		// Austria
		"Vienna": {
			"This city has the oldest zoo in the world, founded in 1752.",
			"Mozart, Beethoven, and Strauss all lived and composed here.",
			"The local coffee house culture is a UNESCO Intangible Cultural Heritage.",
			"The Vienna State Opera performs a different opera every night.",
		},
		// Denmark
		"Copenhagen": {
			"The Tivoli Gardens here inspired Walt Disney to create Disneyland.",
			"This city is considered the most bicycle-friendly in the world.",
			"The Little Mermaid statue has been vandalized over 30 times.",
			"Nyhavn, the colorful harbor, was once a red-light district.",
		},
		// Ireland
		"Dublin": {
			"The Guinness brewery has a 9,000-year lease at only £45 per year.",
			"This city has produced four Nobel Prize winners in literature.",
			"St. Patrick's Cathedral is the largest church in the country.",
			"The Book of Kells, a 1,200-year-old manuscript, is housed here.",
		},
		// Turkey
		"Istanbul": {
			"This is the only city in the world that spans two continents.",
			"The Grand Bazaar has over 4,000 shops and is one of the oldest markets.",
			"The Hagia Sophia has served as a church, mosque, and museum.",
			"Tulips were introduced to Europe from gardens here, not the Netherlands.",
		},
		// UAE
		"Dubai": {
			"The Burj Khalifa is the tallest building in the world at 828 meters.",
			"This city has indoor ski slopes in the middle of the desert.",
			"The Palm Jumeirah added 520 km of beaches to the coastline.",
			"Gold vending machines dispense real gold bars here.",
		},
		// Hong Kong
		"Hong Kong": {
			"This city has the most skyscrapers in the world - over 480 buildings above 150m.",
			"The MTR system is profitable and returns dividends to shareholders.",
			"Dim sum originated here and is traditionally served with tea.",
			"Victoria Peak offers panoramic views and has a tram operating since 1888.",
		},
		// Singapore
		"Singapore": {
			"Chewing gum has been banned here since 1992.",
			"The airport has a butterfly garden with over 1,000 butterflies.",
			"This city-state is one of only three surviving city-states in the world.",
			"The world's first night zoo, the Night Safari, opened here.",
		},
		// Japan
		"Tokyo": {
			"The Shibuya Crossing is the busiest pedestrian intersection in the world.",
			"There are more Michelin-starred restaurants here than in Paris.",
			"Vending machines here sell everything from eggs to umbrellas.",
			"The metro system moves 8.7 million passengers daily.",
		},
		// South Korea
		"Seoul": {
			"The subway system has heated seats in winter.",
			"This city has the fastest internet speeds in the world.",
			"K-pop and Korean Wave (Hallyu) originated from this entertainment hub.",
			"Changdeokgung Palace is a UNESCO World Heritage site.",
		},
		// China
		"Beijing": {
			"The Forbidden City has 9,999 rooms and took 14 years to build.",
			"The Great Wall is about an hour's drive from the city center.",
			"Duck is the signature dish, with restaurants specializing in it for centuries.",
			"Tiananmen Square is the largest public square in the world.",
		},
		"Shanghai": {
			"The Maglev train here reaches 431 km/h, the fastest commercial train.",
			"The Bund features the largest collection of Art Deco buildings outside Europe.",
			"This city has the world's second-tallest building - Shanghai Tower.",
			"The local dialect is unintelligible to Mandarin speakers.",
		},
		// Thailand
		"Bangkok": {
			"The full ceremonial name of this city is 169 letters long.",
			"Floating markets have existed here for over 100 years.",
			"There are over 400 Buddhist temples in the city.",
			"Street food here has earned Michelin stars.",
		},
		// Malaysia
		"Kuala Lumpur": {
			"The Petronas Towers were the tallest buildings in the world from 1998-2004.",
			"Batu Caves nearby feature a 43-meter golden statue.",
			"The city name translates to 'muddy confluence' in Malay.",
			"This is one of the fastest-growing metropolitan regions in Southeast Asia.",
		},
		// India
		"New Delhi": {
			"This city was designed by British architects and completed in 1931.",
			"The Lotus Temple here receives more visitors than the Taj Mahal.",
			"The metro system is the first in India and carries 6 million passengers daily.",
			"India Gate commemorates 70,000 soldiers who died in World War I.",
		},
		"Mumbai": {
			"Bollywood produces more films here annually than Hollywood.",
			"The Dabbawala lunch delivery system has a 99.97% accuracy rate.",
			"This city has the most expensive home in the world - a 27-story private residence.",
			"The Gateway of India was built to welcome King George V in 1924.",
		},
		// Australia
		"Sydney": {
			"The Opera House roof is covered with over 1 million tiles.",
			"The Harbour Bridge is the world's largest steel arch bridge.",
			"Bondi Beach has been a popular destination since the 1850s.",
			"This city hosted the 2000 Summer Olympics.",
		},
		"Melbourne": {
			"This city is considered the coffee capital of Australia.",
			"The laneways are famous for street art and hidden cafes.",
			"Australian Rules Football (AFL) was invented here.",
			"The tram network is the largest urban tram system in the world.",
		},
		// New Zealand
		"Auckland": {
			"More than 50 volcanoes exist within the metropolitan area.",
			"The Sky Tower is the tallest free-standing structure in the Southern Hemisphere.",
			"This city has more boats per capita than any other city in the world.",
			"The America's Cup sailing race has been hosted here multiple times.",
		},
		// Qatar
		"Doha": {
			"The Museum of Islamic Art houses the largest collection of Islamic art.",
			"This city will host the 2022 FIFA World Cup.",
			"The Pearl-Qatar is an artificial island spanning nearly 4 million square meters.",
			"Air conditioning here accounts for 70% of electricity consumption.",
		},
		// Brazil
		"São Paulo": {
			"This is the largest city in the Southern Hemisphere by population.",
			"Japanese immigrants created the largest Japanese community outside Japan here.",
			"Helicopters are commonly used to avoid traffic - there are over 400 registered.",
			"The city has over 12,000 restaurants representing cuisines from around the world.",
		},
		// Argentina
		"Buenos Aires": {
			"Tango was born in the working-class neighborhoods here.",
			"The widest avenue in the world, 9 de Julio, has 16 lanes.",
			"Café Tortoni, opened in 1858, is the city's oldest café.",
			"More psychologists per capita live here than anywhere else in the world.",
		},
		// South Africa
		"Johannesburg": {
			"This city was built on the world's largest gold deposits.",
			"The Apartheid Museum documents South Africa's history of segregation.",
			"This is the largest city not built on a river, lake, or coastline.",
			"The urban forest has over 10 million trees - one of the largest man-made forests.",
		},
		// Egypt
		"Cairo": {
			"The pyramids of Giza are older than this city itself.",
			"The Egyptian Museum houses over 120,000 ancient artifacts.",
			"Traffic here is legendary - honking is a form of communication.",
			"This is the largest city in the Arab world and Africa.",
		},
	}

	// Generic fallback facts for unknown cities
	fallbackFacts = []string{
		"This destination is famous for its unique blend of ancient and modern architecture.",
		"The local cuisine here has been perfected over centuries and attracts food lovers worldwide.",
		"A world-renowned annual festival draws visitors from across the globe to this destination.",
		"The historic old town contains buildings and traditions dating back several hundred years.",
		"This city sits at the crossroads of ancient trade routes that shaped its diverse culture.",
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GetCityFact returns a fact about the specified city
func GetCityFact(city string, seed string) string {
	if city == "" {
		return ""
	}

	// Create cache key
	cacheKey := city + "-" + seed

	// Check cache first
	factCacheMu.RLock()
	if fact, ok := factCache[cacheKey]; ok {
		factCacheMu.RUnlock()
		return fact
	}
	factCacheMu.RUnlock()

	// Get city-specific facts or fallback
	var fact string
	if facts, ok := cityFacts[city]; ok && len(facts) > 0 {
		fact = facts[rand.Intn(len(facts))]
	} else {
		fact = fallbackFacts[rand.Intn(len(fallbackFacts))]
	}

	// Cache the result
	factCacheMu.Lock()
	factCache[cacheKey] = fact
	factCacheMu.Unlock()

	return fact
}

// ClearCache clears the fact cache (useful for testing or memory management)
func ClearCache() {
	factCacheMu.Lock()
	factCache = make(map[string]string)
	factCacheMu.Unlock()
}
