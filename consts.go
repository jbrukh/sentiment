package sentiment

var commonEnglish = []string{
    "to",
    "the",
    "of",
    "and",
    "with",
    "at",
    "a",
    "you",
    "in",
    "it",
    "he",
    "she",
    "it",
    "on",
    "is",
    "for",
    "where", "an", "by", "i", "we", "they", "are", "aren't", "be", "can't", "from", "u",
    "if", "this", "that", "its", "it's", "has", "not", "just", "i'm", "we're", "we", "you're",
    "your", "you'll", "i'll",
}

var twitterTrash = []string{
   "...",
   "rt",
   "&",
   "w/",
   "@",
   "-",
}

func CommonEnglish() []string {
    return commonEnglish
}

func TwitterTrash() []string {
    return twitterTrash
}
