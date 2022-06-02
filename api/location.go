// location.go contains maps and slices of Auckland streets and suburbs.
package api

// street_abbreviations maps common street suffix abbreviations to their uncondensed counterpart.
var street_abbreviations = map[string]string{
	"pl": "Place", "rd": "Road", "st": "Street", "ave": "Avenue", "dr": "Drive",
	"cr": "Crescent", "cres": "Crescent", "cresc": "Crescent", "blvd": "Boulevard",
	"bvd": "Boulevard", "cl": "Close", "hl": "Hill", "tce": "Terrace", "gln": "Glen",
	"pde": "Parade", "ct": "Court", "plza": "Plaza", "prom": "Promenade", "rt": "Retreat",
	"rtt": "Retreat", "rdge": "Ridge", "sq": "Square", "arc": "Arcade", "bwlk": "Boardwalk",
	"ch": "Chase", "cct": "Circuit", "bp": "Bypass", "bypa": "Bypass", "brk": "Break",
	"crst": "Crest", "ent": "Entrance", "esp": "Esplanade", "fwy": "Freeway",
	"glde": "Glade", "gd": "Glade", "gra": "Grange", "ln": "Lane", "av": "Avenue",
}

// suburb_abbreviations maps common suburb suffix abbreviations to their uncondensed counterpart.
var suburb_abbreviations = map[string]string{
	"mt": "Mount", "pt": "Point", "st": "Saint", "cbd": "Central",
}

// suburbs is a slice of most Auckland suburbs.
var suburbs = []string{
	"woodhill forest", "wiri", "windsor park", "whitford", "whenuapai", "wharehine", "whangaripo",
	"whangateau", "whangaparaoa", "weymouth", "westmere", "westlake", "westgate", "west harbour", "westfield",
	"western heights", "wellsford", "wattledowns", "wattle downs", "wattle cove", "waterview", "western springs",
	"warkworth", "waiwera", "waiuku", "waitoki", "waitakere", "wairau valley", "wainui", "waimauku",
	"waimahia landing", "waikowhai", "waiau pa", "waiatarua", "waiake", "wai o taiki bay", "wade heads",
	"unsworth heights", "tuakau", "totara heights", "tuscany estate", "torbay", "totara vale", "tomarata",
	"titirangi", "ti point", "tindalls beach", "three kings", "the gardens", "te papapa", "te hana", "te atatu south",
	"te atatu peninsula", "te atatu", "te arai", "tawharanui peninsula", "taupaki", "tauhoa", "tapora",
	"takapuna", "swanson", "surfdale", "sunnyvale", "sunnynook", "sunnyhills", "stonefields", "stanmore bay",
	"stanley point", "stanley bay", "st johns", "south head", "somerville", "silverdale", "snells beach",
	"silkwood heights", "shelly park", "shelly beach", "shamrock park", "schnapper rock", "sandspit", "saint marys bay",
	"sandringham", "saint lukes", "totara park", "saint johns", "saint heliers", "runciman", "royal oak",
	"royal heights", "rothesay bay", "rosehill", "takanini", "rosedale", "riverhead", "remuera", "redvale",
	"red hill", "red beach", "randwick park", "settlers cove", "pukekohe", "puhoi", "port albert", "ponsonby", "pollok",
	"point wells", "point chevalier", "piha", "ranui", "penrose", "point england", "patumahoe", "parnell", "paremoremo",
	"pinehill", "parau", "parakai", "papatoetoe", "paparimu", "papakura", "panmure", "palm beach", "pakuranga",
	"pakuranga heights", "pakiri", "pahurehure", "paerata", "owairaka", "oteha", "otara", "otahuhu", "ostend",
	"ormiston", "orere point", "oratia", "oranga", "orakei", "opaheke", "onetangi", "one tree hill", "oneroa",
	"onehunga", "omaha", "orewa", "okura", "northpark", "northcross", "northcote", "northcote point", "north harbour",
	"newton", "newmarket", "new windsor", "new lynn", "narrow neck", "muriwai", "mount wellington", "mount roskill",
	"mount albert", "mount eden", "morningside", "mission bay", "millwater", "milford", "murrays bay", "middlemore",
	"murphys heights", "mellons bay", "meadowlands", "meadowbank", "mclaren park", "matakatia", "matakana", "massey",
	"marlborough", "maraetai", "manurewa east", "manurewa", "manukau heights", "manukau central", "manly",
	"manhurangi", "mangere bridge", "mangere east", "mangakura", "makarau", "mairangi bay", "mahurangi west",
	"mahia park", "lynfield", "longford park", "long bay", "lincoln", "leigh", "mahurangi east", "laingholm", "kumeu",
	"konini", "kohimarama", "kingsland", "kingseat", "manukau heads", "kelston", "kawakawa bay", "kaurilands", "kaukapakapa",
	"karaka", "karaka harbourside", "kaipara flats", "hunua", "huntington park", "huia", "karekare", "huapai", "howick",
	"hobsonville", "hillsborough", "hillpark", "hillcrest", "highland park", "highbury", "herald island", "henderson valley",
	"helensville", "hauraki", "hatfields beach", "half moon bay", "gulf harbour", "grey lynn", "herne bay",
	"greenwoods corner", "greenlane", "greenhithe", "green bay", "grafton", "goodwood heights", "greenmeadows", "golflands",
	"glorit", "glenfield", "glendowie", "glendene", "glenbrook", "glen innes", "glen eden", "freemans bay", "flat bush",
	"favona", "farm cove", "fairview heights", "forrest hill", "epsom", "eden valley", "ellerslie", "eden terrace",
	"eastern beach", "east tamaki heights", "drury", "dome forest", "dome valley", "devonport", "dannemora", "east tamaki",
	"dairy flat", "crown hill", "cornwallis", "bethells beach", "cockle bay", "clover park", "clevedon", "clendon park",
	"clarks beach", "cheltenham", "chatswood", "chapel downs", "castor bay", "campbells bay", "burswood", "bucklands beach",
	"buckland", "botany downs", "brookby", "bombay", "blockhouse bay", "blackpool", "birkenhead", "browns bay", "birkdale",
	"big omaha", "belmont", "beachlands", "beach haven", "bayview", "bayswater", "bays water", "balmoral", "awhitu",
	"avondale", "auckland cbd", "army bay", "arkles bay", "ardmore", "botany", "conifer grove", "arch hill", "ararimu",
	"anawhata", "algies bay", "alfriston", "albany", "airport oaks", "henderson", "tamaki", "mangere", "mahurangi",
	"manukau",
}
