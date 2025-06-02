<template>
  <div class="dashboard-wrapper">
    <!-- Sidebar with Products -->
    <div class="sidebar">
      <div class="sidebar-header">
        <h3>🥦 Produkty</h3>
        <div class="search-container">
          <input
              v-model="searchQuery"
              type="text"
              placeholder="Szukaj produktów..."
              class="search-input"
          />
          <span class="search-icon">🔍</span>
        </div>
      </div>

      <div class="products-scroll">
        <div v-for="product in filteredProducts" :key="product.id" class="product-item">
          <div class="product-line">
            <span class="product-name">{{ product.name }}</span>
            <div class="quantity-controls">
              <input
                  type="number"
                  min="0"
                  v-model.number="quantities[product.id]"
                  @input="validate(product.id)"
                  placeholder="g"
                  class="quantity-input"
              />
              <button @click="add(product.id)" class="control-btn add-btn">➕</button>
              <button @click="remove(product.id)" class="control-btn remove-btn">➖</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Panel -->
    <div class="main-panel">
      <!-- Goals Progress Bars -->
      <div class="goals-section">
        <h4>Cele dzienne</h4>
        <div class="goals-bar">
          <div v-for="(goal, key) in goals" :key="key" class="goal-bar">
            <div class="bar-container">
              <div
                  class="bar-fill"
                  :class="key"
                  :style="{ height: barHeight(key) + '%' }"
              ></div>
              <div v-if="barHeight(key) >= 100" class="checkmark">✅</div>
            </div>
            <input
                v-model.number="goals[key]"
                type="number"
                class="goal-input"
            />
            <label class="goal-label">{{ labels[key] }}</label>
            <span class="current-value">{{ Math.round(result?.[key] || 0) }}</span>
          </div>
        </div>
      </div>

      <!-- Chart Section -->
      <div class="chart-section">
        <h4>Makroskładniki</h4>
        <div class="chart-wrapper">
          <div class="pie-chart-container">
            <canvas ref="pieChart" width="200" height="200" class="pie-chart"></canvas>
            <div class="chart-center">
              <div class="total-calories">{{ Math.round(result?.calories || 0) }}</div>
              <div class="calories-label">kcal</div>
            </div>
          </div>

          <div class="legend">
            <div class="legend-item">
              <div class="legend-color protein"></div>
              <span>Białko: {{ Math.round(result?.protein || 0) }}g</span>
            </div>
            <div class="legend-item">
              <div class="legend-color fat"></div>
              <span>Tłuszcz: {{ Math.round(result?.fat || 0) }}g</span>
            </div>
            <div class="legend-item">
              <div class="legend-color carbs"></div>
              <span>Węglowodany: {{ Math.round(result?.carbs || 0) }}g</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, nextTick, watch } from 'vue'
import axios from 'axios'

const products = ref([])
const quantities = reactive({})
const result = ref(null)
const searchQuery = ref('')
const pieChart = ref(null)

const goals = reactive({
  calories: 2000,
  protein: 100,
  fat: 70,
  carbs: 250,
})

const labels = {
  calories: 'Kalorie',
  protein: 'Białko',
  fat: 'Tłuszcz',
  carbs: 'Węglowodany',
}

// Filter products based on search query and show 6 visible, but allow scrolling for more
const filteredProducts = computed(() => {
  const filtered = products.value.filter(product =>
      product.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
  return filtered // Remove slice limit to allow scrolling through all filtered products
})

function validate(id) {
  if (quantities[id] < 0) quantities[id] = 0
  calculate()
}

function add(id) {
  quantities[id] = (quantities[id] || 0) + 50
  calculate()
}

function remove(id) {
  quantities[id] = Math.max(0, (quantities[id] || 0) - 50)
  calculate()
}

function barHeight(key) {
  const val = result.value?.[key] || 0
  return Math.min(100, (val / goals[key]) * 100)
}

function drawPieChart() {
  if (!pieChart.value || !result.value) return

  const canvas = pieChart.value
  const ctx = canvas.getContext('2d')
  const centerX = canvas.width / 2
  const centerY = canvas.height / 2
  const radius = 80

  ctx.clearRect(0, 0, canvas.width, canvas.height)

  const protein = result.value.protein || 0
  const fat = result.value.fat || 0
  const carbs = result.value.carbs || 0
  const total = protein + fat + carbs

  if (total === 0) return

  const colors = {
    protein: '#a0c4ff',
    fat: '#ffc6a0',
    carbs: '#98e8a0'
  }

  let currentAngle = -Math.PI / 2

  // Draw protein slice
  const proteinAngle = (protein / total) * 2 * Math.PI
  ctx.beginPath()
  ctx.moveTo(centerX, centerY)
  ctx.arc(centerX, centerY, radius, currentAngle, currentAngle + proteinAngle)
  ctx.closePath()
  ctx.fillStyle = colors.protein
  ctx.fill()

  // Draw fat slice
  currentAngle += proteinAngle
  const fatAngle = (fat / total) * 2 * Math.PI
  ctx.beginPath()
  ctx.moveTo(centerX, centerY)
  ctx.arc(centerX, centerY, radius, currentAngle, currentAngle + fatAngle)
  ctx.closePath()
  ctx.fillStyle = colors.fat
  ctx.fill()

  // Draw carbs slice
  currentAngle += fatAngle
  const carbsAngle = (carbs / total) * 2 * Math.PI
  ctx.beginPath()
  ctx.moveTo(centerX, centerY)
  ctx.arc(centerX, centerY, radius, currentAngle, currentAngle + carbsAngle)
  ctx.closePath()
  ctx.fillStyle = colors.carbs
  ctx.fill()

  // Draw inner circle to create donut effect
  ctx.beginPath()
  ctx.arc(centerX, centerY, 45, 0, 2 * Math.PI)
  ctx.fillStyle = 'rgba(255, 255, 255, 0.9)'
  ctx.fill()
}

async function calculate() {
  const items = Object.entries(quantities)
      .filter(([_, grams]) => grams > 0)
      .map(([id, grams]) => ({ product_id: Number(id), grams }))

  if (!items.length) {
    result.value = { calories: 0, protein: 0, fat: 0, carbs: 0 }
    drawPieChart()
    return
  }

  try {
    const res = await axios.post('/calculate', { items })
    result.value = res.data
    await nextTick()
    drawPieChart()
  } catch (error) {
    console.error('Error calculating nutrition:', error)
  }
}

watch(result, () => {
  nextTick(() => drawPieChart())
}, { deep: true })

onMounted(async () => {
  try {
    const res = await axios.get('/products')
    products.value = res.data
    calculate()
  } catch (error) {
    console.error('Error loading products:', error)
  }
})
</script>

<style scoped>
.dashboard-wrapper {
  display: flex;
  height: 100vh;
  max-height: 100vh;
  width: 100%;
  max-width: 100%;
  padding: 1.5rem;
  gap: 1.5rem;
  background: rgba(139, 231, 139, 0.15);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

/* Sidebar Styles - WITH SCROLLING RESTORED */
.sidebar {
  flex: 1;
  height: 100%;
  max-height: 100%;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.sidebar-header {
  flex-shrink: 0;
  margin-bottom: 1rem;
}

.sidebar-header h3 {
  margin: 0 0 1rem 0;
  color: #2d5016;
  font-size: 1.4rem;
  font-weight: 700;
}

.search-container {
  position: relative;
}

.search-input {
  width: 100%;
  padding: 0.75rem 2.5rem 0.75rem 1rem;
  border: 2px solid #e8f5e8;
  border-radius: 12px;
  font-size: 0.9rem;
  background: rgba(255, 255, 255, 0.9);
  transition: all 0.3s ease;
}

.search-input:focus {
  outline: none;
  border-color: #8de78d;
  background: white;
  box-shadow: 0 0 0 3px rgba(141, 231, 141, 0.2);
}

.search-icon {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  color: #666;
}

/* SCROLLING RESTORED - shows about 6 products, scroll for more */
.products-scroll {
  flex: 1;
  overflow-y: auto;
  padding-right: 0.5rem;
  min-height: 0;
  max-height: calc(100% - 120px); /* Reserve space for header */
}

.products-scroll::-webkit-scrollbar {
  width: 6px;
}

.products-scroll::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.products-scroll::-webkit-scrollbar-thumb {
  background: #8de78d;
  border-radius: 3px;
}

.product-item {
  margin-bottom: 0.8rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 12px;
  border: 1px solid #e8f5e8;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.product-item:hover {
  background: rgba(255, 255, 255, 0.95);
  border-color: #8de78d;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.product-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.product-name {
  flex: 1;
  font-weight: 500;
  color: #2d5016;
  font-size: 0.9rem;
}

.quantity-controls {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.quantity-input {
  width: 55px;
  padding: 0.4rem;
  border: 1px solid #ddd;
  border-radius: 8px;
  text-align: center;
  font-size: 0.85rem;
}

.control-btn {
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.8rem;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.add-btn {
  background: #8de78d;
  color: white;
}

.add-btn:hover {
  background: #6ed46e;
  transform: scale(1.1);
}

.remove-btn {
  background: #ff9999;
  color: white;
}

.remove-btn:hover {
  background: #ff7777;
  transform: scale(1.1);
}

/* Main Panel Styles */
.main-panel {
  flex: 2;
  height: 100%;
  max-height: 100%;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  overflow: hidden;
}

.goals-section, .chart-section {
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  flex-shrink: 0;
}

.goals-section h4, .chart-section h4 {
  margin: 0 0 1.2rem 0;
  color: #2d5016;
  font-size: 1.2rem;
  font-weight: 600;
  text-align: center;
}

.goals-bar {
  display: flex;
  justify-content: center;
  gap: 1.5rem;
}

.goal-bar {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.4rem;
}

.bar-container {
  width: 35px;
  height: 100px;
  background: #f0f8f0;
  border-radius: 18px;
  position: relative;
  border: 2px solid #e8f5e8;
  overflow: hidden;
}

.bar-fill {
  position: absolute;
  bottom: 0;
  width: 100%;
  border-radius: 16px;
  transition: all 0.5s ease;
}

.bar-fill.calories { background: linear-gradient(to top, #ff6b6b, #ffa8a8); }
.bar-fill.protein { background: linear-gradient(to top, #a0c4ff, #c4d9ff); }
.bar-fill.fat { background: linear-gradient(to top, #ffc6a0, #ffe0c4); }
.bar-fill.carbs { background: linear-gradient(to top, #98e8a0, #c4f0c8); }

.checkmark {
  position: absolute;
  top: -22px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 1rem;
}

.goal-input {
  width: 65px;
  padding: 0.25rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  text-align: center;
  font-size: 0.85rem;
}

.goal-label {
  font-size: 0.8rem;
  font-weight: 500;
  color: #2d5016;
}

.current-value {
  font-size: 0.75rem;
  color: #666;
  background: #f0f8f0;
  padding: 0.15rem 0.4rem;
  border-radius: 10px;
}

/* Chart Styles */
.chart-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.2rem;
}

.pie-chart-container {
  position: relative;
  display: inline-block;
}

.pie-chart {
  border-radius: 50%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.chart-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
}

.total-calories {
  font-size: 1.6rem;
  font-weight: 700;
  color: #2d5016;
}

.calories-label {
  font-size: 0.8rem;
  color: #666;
  margin-top: -0.2rem;
}

.legend {
  display: flex;
  gap: 1.2rem;
  justify-content: center;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.85rem;
  color: #2d5016;
}

.legend-color {
  width: 14px;
  height: 14px;
  border-radius: 50%;
}

.legend-color.protein { background: #a0c4ff; }
.legend-color.fat { background: #ffc6a0; }
.legend-color.carbs { background: #98e8a0; }
</style>