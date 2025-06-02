<template>
  <div class="about-wrapper">
    <div class="about-container">
      <!-- Header z animowanym logo -->
      <div class="about-header">
        <div class="logo-section" @click="triggerEasterEgg">
          <div class="animated-logo" :class="{ bounce: logoBouncied }">
            <span class="logo-text">🏋️‍♂️ trAIning buddy</span>
            <div class="ai-highlight">AI</div>
          </div>
          <p class="tagline">{{ currentTagline }}</p>
        </div>
      </div>

      <!-- Licznik odwiedzin -->
      <div class="visit-counter" @click="incrementVisits">
        👁️ Odwiedziny: {{ visitCount }} | Kliknięcia logo: {{ logoClicks }}
      </div>

      <!-- Sekcja Misji -->
      <section class="mission-section">
        <div class="section-header">
          <h2>🎯 Nasza Misja</h2>
          <button @click="changeMission" class="fun-btn">🎲 Losuj misję</button>
        </div>
        <div class="mission-card">
          <p>{{ currentMission }}</p>
          <div class="mission-stats">
            <span>📊 Skuteczność: {{ missionStats.effectiveness }}%</span>
            <span>🚀 Motywacja: {{ missionStats.motivation }}%</span>
            <span>😂 Poziom żartów: {{ missionStats.jokes }}%</span>
          </div>
        </div>
      </section>

      <!-- Sekcja Zespołu -->
      <section class="team-section">
        <h2>👥 Nasz Wspaniały Zespół</h2>
        <div class="team-grid">
          <div
              v-for="(member, index) in teamMembers"
              :key="index"
              class="team-card"
              @click="rotateMember(index)"
              :style="{ transform: member.rotated ? 'rotateY(180deg)' : 'rotateY(0deg)' }"
          >
            <div class="card-front">
              <div class="avatar" :style="{ backgroundColor: member.color }">
                {{ member.emoji }}
              </div>
              <h3>{{ member.name }}</h3>
              <p class="role">{{ member.role }}</p>
              <p class="motto">{{ member.motto }}</p>
            </div>
            <div class="card-back">
              <h3>Tajne info! 🕵️</h3>
              <p>{{ member.secret }}</p>
              <div class="skills">
                <span v-for="skill in member.skills" :key="skill" class="skill-tag">
                  {{ skill }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Sekcja Ciekawostek -->
      <section class="facts-section">
        <h2>🤓 Niesamowite Fakty</h2>
        <div class="facts-container">
          <div
              v-for="(fact, index) in facts"
              :key="index"
              class="fact-card"
              @click="toggleFact(index)"
              :class="{ revealed: fact.revealed }"
          >
            <div class="fact-icon">{{ fact.icon }}</div>
            <div class="fact-content">
              <h4>{{ fact.title }}</h4>
              <p v-if="fact.revealed">{{ fact.content }}</p>
              <p v-else class="click-to-reveal">Kliknij aby odkryć prawdę!</p>
            </div>
          </div>
        </div>
      </section>

      <!-- Generator Motywacyjny -->
      <section class="motivation-section">
        <h2>💪 Generator Motywacji</h2>
        <div class="motivation-generator">
          <div class="current-quote">
            <p>"{{ currentQuote }}"</p>
            <span class="quote-author">- {{ currentAuthor }}</span>
          </div>
          <button @click="generateQuote" class="generate-btn">
            🎯 Nowa motywacja!
          </button>
          <div class="quote-stats">
            Wygenerowano już {{ quotesGenerated }} cytatów motywacyjnych!
          </div>
        </div>
      </section>

      <!-- Easter Egg Section (uproszczona) -->
      <section v-if="easterEggActivated" class="easter-egg-section">
        <h2>🥚 Gratulacje! Znalazłeś Easter Egg!</h2>
        <div class="easter-content">
          <div class="dancing-emoji">🎉</div>
          <p>Wspaniale! Udało Ci się odkryć nasz sekretny easter egg!</p>
        </div>
      </section>

      <!-- Kontakt -->
      <section class="contact-section">
        <h2>📞 Kontakt</h2>
        <div class="contact-grid">
          <div class="contact-card">
            <div class="contact-icon">📧</div>
            <h4>Email</h4>
            <p>kontakt@training-buddy.pl</p>
            <p class="fun-note">(Odpowiadamy w ciągu 2-3 dni roboczych, chyba że akurat trenujemy)</p>
          </div>
          <div class="contact-card">
            <div class="contact-icon">📱</div>
            <h4>Telefon</h4>
            <p>+48 123 456 789</p>
            <p class="fun-note">(Najlepiej dzwonić po 16:00, wcześniej jesteśmy na siłowni)</p>
          </div>
          <div class="contact-card">
            <div class="contact-icon">📍</div>
            <h4>Adres</h4>
            <p>ul. Mięśniowa 42</p>
            <p>00-001 Warszawa</p>
            <p class="fun-note">(Tuż obok najlepszej kebabowej w mieście)</p>
          </div>
        </div>
      </section>

      <!-- Footer -->
      <footer class="about-footer">
        <div class="footer-content">
          <p class="footer-main">© 2024 trAIning buddy - Gdzie AI spotyka się z gainzami! 💪</p>
          <div class="footer-stats">
            <span>🏋️ {{ Math.floor(Math.random() * 1000) + 500 }} zadowolonych użytkowników</span>
            <span>📈 {{ Math.floor(Math.random() * 50) + 25 }}% więcej motywacji</span>
            <span>☕ {{ Math.floor(Math.random() * 100) + 200 }} wypitych kaw podczas tworzenia</span>
          </div>
        </div>
      </footer>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

// Reactive data
const visitCount = ref(0)
const logoClicks = ref(0)
const logoBouncied = ref(false)
const quotesGenerated = ref(0)
const easterEggActivated = ref(false)

// Taglines
const taglines = [
  "Gdzie sztuczna inteligencja spotyka się z prawdziwą siłą!",
  "AI + Mięśnie = Sukces!",
  "Trenuj mądrze, trenuj z AI!",
  "Twój cyfrowy personal trainer!",
  "Gdzie technologia napędza transformację!"
]
const currentTagline = ref(taglines[0])

// Missions
const missions = [
  "Pomagamy ludziom osiągać ich cele fitness poprzez połączenie najnowszych technologii AI z sprawdzonymi metodami treningowymi. A przy okazji staramy się nie umrzeć ze śmiechu podczas projektowania interfejsów.",
  "Naszą misją jest stworzenie świata, w którym każdy ma dostęp do inteligentnego trenera personalnego, który nigdy nie ma złego dnia i nie zjada Twojego jedzenia z lodówki.",
  "Wierzymy, że fitness powinien być przyjemny, skuteczny i dostępny dla wszystkich. Dlatego stworzyliśmy AI, które Cię motywuje zamiast osądzać za zjedzenie całej pizzy wczoraj.",
  "Łączymy nauki o sporcie z technologią przyszłości, aby pomóc Ci wyrzeźbić ciało marzeń. Ostrzeżenie: efekty uboczne mogą obejmować zwiększoną pewność siebie."
]
const currentMission = ref(missions[0])

// Mission stats (zawsze śmieszne i losowe)
const missionStats = ref({
  effectiveness: Math.floor(Math.random() * 20) + 85,
  motivation: Math.floor(Math.random() * 15) + 90,
  jokes: Math.floor(Math.random() * 30) + 70
})

// Team members
const teamMembers = ref([
  {
    name: "Dr. AI Gainz",
    role: "Chief AI Officer",
    emoji: "🤖",
    color: "#FF6B6B",
    motto: "Beep boop, let's get swole!",
    secret: "Lubi słuchać heavy metalu podczas kodowania algorytmów",
    skills: ["Machine Learning", "Deadlifts", "Debugging"],
    rotated: false
  },
  {
    name: "Flex Rodriguez",
    role: "Head of Motivation",
    emoji: "💪",
    color: "#4ECDC4",
    motto: "No pain, no gain, no brain!",
    secret: "Jego prawdziwe imię to Stanisław i pochodzi z Końskich",
    skills: ["Krzyczenie", "Protein shakes", "Inspiracja"],
    rotated: false
  },
  {
    name: "Sarah Squats",
    role: "UX/UI Designer",
    emoji: "🎨",
    color: "#45B7D1",
    motto: "Design so good, it's almost illegal",
    secret: "Projektuje wszystkie interfejsy w pozycji squat",
    skills: ["Figma", "Photoshop", "Bulgarian Split Squats"],
    rotated: false
  },
  {
    name: "Code McBiceps",
    role: "Full Stack Developer",
    emoji: "👨‍💻",
    color: "#F7DC6F",
    motto: "Lifting weights and raising exceptions",
    secret: "Kiedyś napisał cały backend jedną ręką podczas robienia bicepsów",
    skills: ["Vue.js", "Go", "Biceps curls"],
    rotated: false
  },
  {
    name: "Cardio Queen",
    role: "Marketing Manager",
    emoji: "👸",
    color: "#BB8FCE",
    motto: "Marketing so fast, it's practically cardio",
    secret: "Może biegać i tworzyć kampanie marketingowe jednocześnie",
    skills: ["Social Media", "Marathon", "Lead Generation"],
    rotated: false
  },
  {
    name: "Protein Pete",
    role: "Nutritionist",
    emoji: "🥤",
    color: "#82E0AA",
    motto: "You are what you eat, so I'm basically protein",
    secret: "Zna na pamięć wartości odżywcze każdego jedzenia na świecie",
    skills: ["Nutrition", "Meal Prep", "Whey Protein Mixing"],
    rotated: false
  }
])

// Facts
const facts = ref([
  {
    title: "Fakty o AI",
    content: "Nasze AI trenowało się na ponad 10 milionach zdjęć mięśni. Teraz ma kompleksy.",
    icon: "🤖",
    revealed: false
  },
  {
    title: "Sekret sukcesu",
    content: "73% naszego kodu zostało napisane pod wpływem przedtreningówki. Nie polecamy w domu.",
    icon: "⚡",
    revealed: false
  },
  {
    title: "Statystyki użytkowników",
    content: "Średni użytkownik trAIning buddy podnosi o 23% więcej niż przeciętny człowiek. Może to być przypadek.",
    icon: "📊",
    revealed: false
  },
  {
    title: "Easter Egg",
    content: "W aplikacji ukryliśmy 17 easter eggów. Ten jest 18-ty, nieplanowany.",
    icon: "🥚",
    revealed: false
  },
  {
    title: "Prawda o teamie",
    content: "Nasz zespół spożywa łącznie 847 gram białka dziennie. To około 2 kurczaków.",
    icon: "🐔",
    revealed: false
  },
  {
    title: "Techniczne detale",
    content: "Aplikacja została napisana w Vue.js, ponieważ React nie umiał zrobić squata.",
    icon: "⚙️",
    revealed: false
  }
])

// Quotes
const motivationalQuotes = [
  { text: "Twoje ciało potrafi więcej niż myślisz. Twój umysł też, ale on akurat jest zajęty scrollowaniem TikToka", author: "Guru Fitness 2024" },
  { text: "Jedyny zły trening to ten, którego nie robiłeś", author: "Captain Obvious" },
  { text: "Nie porównuj się z innymi. Porównaj się z wczorajszym sobą. On był słabszy", author: "Mirror Mirror" },
  { text: "Sukces to 1% inspiracji i 99% motywacji do wstawania z łóżka", author: "Edison's Trainer" },
  { text: "Jeśli nie spocisz się podczas treningu, prawdopodobnie robisz yoga", author: "Yoga Hater" },
  { text: "Ból przemija, ale screenshoty z aplikacji fitness zostają na zawsze", author: "Social Media Guru" },
  { text: "Nie ma shortcuts do miejsca, które warto odwiedzić. Chyba że masz windę", author: "Practical Phil" }
]

const currentQuote = ref(motivationalQuotes[0].text)
const currentAuthor = ref(motivationalQuotes[0].author)

// Methods
const incrementVisits = () => {
  visitCount.value++
  localStorage.setItem('training-buddy-visits', visitCount.value.toString())
}

const triggerEasterEgg = () => {
  logoClicks.value++
  logoBouncied.value = true
  setTimeout(() => logoBouncied.value = false, 600)

  if (logoClicks.value >= 5) {
    easterEggActivated.value = true
  }

  // Rotuj tagline co 3 kliknięcia
  if (logoClicks.value % 3 === 0) {
    const randomIndex = Math.floor(Math.random() * taglines.length)
    currentTagline.value = taglines[randomIndex]
  }
}

const changeMission = () => {
  const randomIndex = Math.floor(Math.random() * missions.length)
  currentMission.value = missions[randomIndex]

  // Regeneruj statystyki
  missionStats.value = {
    effectiveness: Math.floor(Math.random() * 20) + 85,
    motivation: Math.floor(Math.random() * 15) + 90,
    jokes: Math.floor(Math.random() * 30) + 70
  }
}

const rotateMember = (index) => {
  teamMembers.value[index].rotated = !teamMembers.value[index].rotated
}

const toggleFact = (index) => {
  facts.value[index].revealed = !facts.value[index].revealed
}

const generateQuote = () => {
  const randomIndex = Math.floor(Math.random() * motivationalQuotes.length)
  currentQuote.value = motivationalQuotes[randomIndex].text
  currentAuthor.value = motivationalQuotes[randomIndex].author
  quotesGenerated.value++
}

// Lifecycle
onMounted(() => {
  // Ładuj liczbę odwiedzin z localStorage
  const savedVisits = localStorage.getItem('training-buddy-visits')
  if (savedVisits) {
    visitCount.value = parseInt(savedVisits)
  }

  // Automatycznie zwiększ liczbę odwiedzin
  incrementVisits()

  // Losuj początkowy cytat
  generateQuote()

  // Losuj początkowy tagline
  const randomIndex = Math.floor(Math.random() * taglines.length)
  currentTagline.value = taglines[randomIndex]
})
</script>

<style scoped>
.about-wrapper {
  padding: 2rem;
  background: rgba(139, 231, 139, 0.1);
  backdrop-filter: blur(10px);
  min-height: 100vh;
  width: 100%;
}

.about-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Header */
.about-header {
  text-align: center;
  padding: 2rem 0;
  background: rgba(255, 255, 255, 0.65);
  border-radius: 20px;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.logo-section {
  cursor: pointer;
  transition: transform 0.3s ease;
}

.logo-section:hover {
  transform: scale(1.05);
}

.animated-logo {
  position: relative;
  display: inline-block;
  transition: transform 0.6s ease;
}

.animated-logo.bounce {
  animation: bounce 0.6s ease;
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-30px);
  }
  60% {
    transform: translateY(-15px);
  }
}

.logo-text {
  font-size: 3rem;
  font-weight: bold;
  background: linear-gradient(45deg, #FF6B6B, #4ECDC4, #45B7D1);
  background-size: 200% 200%;
  animation: gradientShift 3s ease infinite;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

@keyframes gradientShift {
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
}

.ai-highlight {
  position: absolute;
  top: -10px;
  right: -20px;
  background: #FF6B6B;
  color: white;
  padding: 0.25rem 0.5rem;
  border-radius: 10px;
  font-size: 0.8rem;
  font-weight: bold;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}

.tagline {
  font-size: 1.2rem;
  color: #2d5016;
  font-weight: 600;
  margin-top: 1rem;
  font-style: italic;
  text-shadow: 0 1px 2px rgba(255, 255, 255, 0.5);

}

/* Visit Counter */
.visit-counter {
  background: rgba(255, 255, 255, 0.95);
  padding: 1rem;
  border-radius: 15px;
  text-align: center;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  color: #2d5016;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.visit-counter:hover {
  background: rgba(255, 255, 255, 1);
  border-color: #4ECDC4;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
}

/* Sections */
section {
  background: rgba(255, 255, 255, 0.95);
  padding: 2rem;
  border-radius: 20px;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

h2 {
  font-size: 2rem;
  color: #2d5016;
  margin: 0;
  font-weight: bold;
}

.fun-btn, .generate-btn, .secret-btn {
  background: linear-gradient(45deg, #FF6B6B, #4ECDC4);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 25px;
  cursor: pointer;
  font-weight: bold;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
}

.fun-btn:hover, .generate-btn:hover, .secret-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}

/* Mission Section */
.mission-card {
  background: rgba(245, 245, 245, 0.9);
  padding: 2rem;
  border-radius: 15px;
  margin-top: 1rem;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.mission-card p {
  font-size: 1.1rem;
  line-height: 1.6;
  margin-bottom: 1.5rem;
  color: #333;
}

.mission-stats {
  display: flex;
  gap: 2rem;
  flex-wrap: wrap;
}

.mission-stats span {
  background: rgba(78, 205, 196, 0.3);
  color: #2d5016;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-weight: bold;
  border: 1px solid rgba(78, 205, 196, 0.5);
}

/* Team Section */
.team-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}

.team-card {
  perspective: 1000px;
  height: 300px;
  cursor: pointer;
}

.card-front, .card-back {
  position: absolute;
  width: 100%;
  height: 100%;
  backface-visibility: hidden;
  border-radius: 15px;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  transition: transform 0.6s ease;
}

.card-front {
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.card-front h3 {
  color: #2d5016;
  margin: 0.5rem 0;
}

.card-back {
  background: rgba(69, 183, 209, 0.95);
  color: white;
  transform: rotateY(180deg);
  border: 1px solid rgba(69, 183, 209, 0.3);
}

.team-card:hover .card-front {
  transform: rotateY(180deg);
}

.team-card:hover .card-back {
  transform: rotateY(0deg);
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  margin-bottom: 1rem;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
}

.role {
  color: #666;
  font-style: italic;
  margin: 0.5rem 0;
}

.motto {
  font-size: 0.9rem;
  color: #888;
  font-style: italic;
}

.skills {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1rem;
}

.skill-tag {
  background: rgba(255, 255, 255, 0.3);
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
  font-size: 0.8rem;
}

/* Facts Section */
.facts-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-top: 2rem;
}

.fact-card {
  background: rgba(245, 245, 245, 0.9);
  padding: 1.5rem;
  border-radius: 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.fact-card:hover {
  transform: translateY(-5px);
  border-color: #4ECDC4;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

.fact-card.revealed {
  background: rgba(78, 205, 196, 0.2);
  border-color: #4ECDC4;
}

.fact-icon {
  font-size: 2rem;
  flex-shrink: 0;
}

.fact-content h4 {
  margin: 0 0 0.5rem 0;
  color: #2d5016;
}

.fact-content p {
  color: #333;
  margin: 0;
}

.click-to-reveal {
  color: #666 !important;
  font-style: italic;
}

/* Motivation Section */
.motivation-generator {
  text-align: center;
}

.current-quote {
  background: rgba(245, 245, 245, 0.9);
  padding: 2rem;
  border-radius: 15px;
  margin-bottom: 2rem;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.current-quote p {
  font-size: 1.3rem;
  font-style: italic;
  margin-bottom: 1rem;
  color: #2d5016;
}

.quote-author {
  color: #666;
  font-weight: bold;
}

.quote-stats {
  margin-top: 1rem;
  color: #666;
  font-style: italic;
}

/* Easter Egg Section */
.easter-egg-section {
  background: linear-gradient(45deg, #FF6B6B, #4ECDC4, #45B7D1) !important;
  color: white;
  text-align: center;
  animation: rainbow 3s ease infinite;
}

.easter-egg-section h2 {
  color: white !important;
}

@keyframes rainbow {
  0% { filter: hue-rotate(0deg); }
  100% { filter: hue-rotate(360deg); }
}

.dancing-emoji {
  font-size: 4rem;
  animation: dance 1s ease infinite;
  margin: 1rem 0;
}

@keyframes dance {
  0%, 100% { transform: rotate(0deg) scale(1); }
  25% { transform: rotate(-10deg) scale(1.1); }
  75% { transform: rotate(10deg) scale(1.1); }
}

.easter-content p {
  color: white !important;
  font-size: 1.2rem;
}

/* Contact Section */
.contact-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}

.contact-card {
  background: rgba(245, 245, 245, 0.9);
  padding: 2rem;
  border-radius: 15px;
  text-align: center;
  transition: all 0.3s ease;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.contact-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
}

.contact-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.contact-card h4 {
  margin: 0 0 1rem 0;
  color: #2d5016;
}

.contact-card p {
  color: #333;
  margin: 0.25rem 0;
}

.fun-note {
  font-size: 0.9rem;
  color: #666 !important;
  font-style: italic;
  margin-top: 0.5rem;
}

/* Footer */
.about-footer {
  background: rgba(255, 255, 255, 0.95) !important;
  color: #333 !important;
  margin-top: 2rem;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 20px;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
  padding: 2rem 1.5rem;
  min-height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.footer-content {
  text-align: center;
}

.footer-main {
  color: #2d5016 !important;
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
}

.footer-stats {
  display: flex;
  justify-content: center;
  gap: 2rem;
  margin-top: 1rem;
  flex-wrap: wrap;
}

.footer-stats span {
  background: rgba(78, 205, 196, 0.15);
  color: #2d5016 !important;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 500;
  border: 1px solid rgba(78, 205, 196, 0.3);
  transition: all 0.3s ease;
}

.footer-stats span:hover {
  background: rgba(78, 205, 196, 0.25);
  transform: translateY(-2px);
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.1);
}

/* Responsive Design */
@media (max-width: 768px) {
  .about-wrapper {
    padding: 1rem;
  }

  .logo-text {
    font-size: 2rem;
  }

  .mission-stats {
    flex-direction: column;
    gap: 1rem;
  }

  .team-grid {
    grid-template-columns: 1fr;
  }

  .facts-container {
    grid-template-columns: 1fr;
  }

  .contact-grid {
    grid-template-columns: 1fr;
  }

  .footer-stats {
    flex-direction: column;
    gap: 1rem;
  }

  .section-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
}
</style>