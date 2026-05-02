<template>
  <v-container class="py-8">
    <v-row justify="center">
      <v-col cols="12" xl="10">
        <v-card class="overflow-hidden">
          <v-tabs v-model="tab" bg-color="white" color="primary" border-b>
            <v-tab value="theory" class="text-none font-weight-bold">Теория</v-tab>
            <v-tab value="practice" class="text-none font-weight-bold">Задания</v-tab>
          </v-tabs>

          <v-card-text class="pa-6">
            <v-window v-model="tab">
              <v-window-item value="theory">
                <div ref="theoryRoot">
                  <v-card
                    v-for="section in theorySections"
                    :key="section.id"
                    variant="flat"
                    class="pa-6 mb-8 theory-card"
                  >
                    <div class="d-flex align-center mb-4">
                      <v-icon :color="section.accentColor" size="30" class="mr-3">{{ section.icon }}</v-icon>
                      <h2 class="text-h5 text-md-h4 font-weight-bold text-text-primary">{{ section.title }}</h2>
                    </div>

                    <template v-if="section.layout === 'stack'">
                      <div class="text-body-1 leading-relaxed text-text-primary">
                        <div
                          v-for="(p, idx) in section.paragraphsHtml"
                          :key="`${section.id}-p-${idx}`"
                          class="mb-3"
                          v-html="p"
                        ></div>
                      </div>

                      <v-sheet
                        v-if="section.formulaLatex"
                        class="formula-block mt-4 mb-2"
                        rounded="lg"
                        border
                      >
                        <div class="math-font text-body-1" v-html="section.formulaLatex"></div>
                      </v-sheet>

                      <div class="illustration-card mt-5">
                        <component :is="section.illustrationComponent" />
                      </div>

                      <v-expansion-panels
                        v-if="section.notes?.length"
                        multiple
                        variant="accordion"
                        class="mt-5 notes-panels"
                        @update:modelValue="renderTheoryMath"
                      >
                        <v-expansion-panel
                          v-for="(note, noteIdx) in section.notes"
                          :key="`${section.id}-note-${noteIdx}`"
                          class="note-panel"
                        >
                          <v-expansion-panel-title class="text-subtitle-2 font-weight-bold text-text-primary">
                            <v-icon :color="section.accentColor" size="18" class="mr-2">mdi-note-text-outline</v-icon>
                            {{ note.title }}
                          </v-expansion-panel-title>
                          <v-expansion-panel-text>
                            <div class="text-body-2 leading-relaxed text-text-primary" v-html="note.bodyHtml"></div>
                          </v-expansion-panel-text>
                        </v-expansion-panel>
                      </v-expansion-panels>
                    </template>

                    <v-row v-else class="align-start" dense>
                      <v-col cols="12" md="7">
                        <div class="text-body-1 leading-relaxed text-text-primary">
                          <div
                            v-for="(p, idx) in section.paragraphsHtml"
                            :key="`${section.id}-p-${idx}`"
                            class="mb-3"
                            v-html="p"
                          ></div>
                        </div>

                        <v-sheet
                          v-if="section.formulaLatex"
                          class="formula-block mt-4 mb-2"
                          rounded="lg"
                          border
                        >
                          <div class="math-font text-body-1" v-html="section.formulaLatex"></div>
                        </v-sheet>

                        <v-expansion-panels
                          v-if="section.notes?.length"
                          multiple
                          variant="accordion"
                          class="mt-4 notes-panels"
                          @update:modelValue="renderTheoryMath"
                        >
                          <v-expansion-panel
                            v-for="(note, noteIdx) in section.notes"
                            :key="`${section.id}-note-${noteIdx}`"
                            class="note-panel"
                          >
                            <v-expansion-panel-title class="text-subtitle-2 font-weight-bold text-text-primary">
                              <v-icon :color="section.accentColor" size="18" class="mr-2">mdi-note-text-outline</v-icon>
                              {{ note.title }}
                            </v-expansion-panel-title>
                            <v-expansion-panel-text>
                              <div class="text-body-2 leading-relaxed text-text-primary" v-html="note.bodyHtml"></div>
                            </v-expansion-panel-text>
                          </v-expansion-panel>
                        </v-expansion-panels>
                      </v-col>

                      <v-col cols="12" md="5" class="mt-5 mt-md-0">
                        <div class="illustration-card">
                          <component :is="section.illustrationComponent" />
                        </div>
                      </v-col>
                    </v-row>
                  </v-card>
                </div>
              </v-window-item>

              <v-window-item value="practice">
                <div v-for="task in tasks" :key="task.id" class="mb-10">
                  <v-card variant="flat" class="pa-6 border-thin bg-surface-light rounded-xl">
                    <div class="d-flex justify-space-between align-start mb-4">
                      <h3 class="text-h5 font-weight-bold">{{ task.title }}</h3>
                      <v-chip
                        v-if="results[task.id] !== undefined"
                        :color="results[task.id] ? 'success' : 'error'"
                        size="small"
                        variant="flat"
                        class="font-weight-bold"
                      >
                        {{ results[task.id] ? 'Верно' : 'Ошибка' }}
                      </v-chip>
                    </div>
                    
                    <p class="text-body-1 mb-6 text-text-secondary">{{ task.description }}</p>
                    
                    <v-row align="center">
                      <v-col cols="12" sm="8" md="6">
                        <v-text-field
                          v-model.number="userAnswers[task.id]"
                          label="Ваш ответ"
                          type="number"
                          variant="outlined"
                          density="comfortable"
                          class="math-font"
                          :color="getAnswerColor(task.id)"
                          :error="results[task.id] === false"
                          :success="results[task.id] === true"
                          persistent-placeholder
                          placeholder="0.00"
                          hide-details
                          @keyup.enter="checkAnswer(task.id)"
                        >
                          <template v-slot:append-inner>
                            <v-icon :color="getAnswerColor(task.id)">{{ getAnswerIcon(task.id) }}</v-icon>
                          </template>
                        </v-text-field>
                      </v-col>
                      <v-col cols="12" sm="4" md="3">
                        <v-btn
                          block
                          color="primary"
                          height="48"
                          class="font-weight-bold text-none"
                          elevation="0"
                          @click="checkAnswer(task.id)"
                        >Проверить</v-btn>
                      </v-col>
                      <v-col cols="12" md="1" class="text-center">
                        <v-tooltip location="top">
                          <template v-slot:activator="{ props }">
                            <v-icon v-bind="props" color="secondary" size="24">mdi-help-circle-outline</v-icon>
                          </template>
                          <span>Требуемая точность: ±{{ task.precision || 0.01 }}</span>
                        </v-tooltip>
                      </v-col>
                    </v-row>
                    
                    <v-expand-transition>
                      <div v-if="results[task.id] === true" class="mt-6">
                        <div class="text-subtitle-2 font-weight-bold mb-3 text-success d-flex align-center">
                          <v-icon size="18" class="mr-2">mdi-chart-line</v-icon>
                          Визуализация решения
                        </div>
                        <v-card variant="outlined" class="pa-2 bg-white" style="min-height: 320px;">
                          <div v-if="taskLoading[task.id]" class="d-flex justify-center align-center h-100">
                            <v-progress-circular indeterminate color="primary"></v-progress-circular>
                          </div>
                          <ChartView
                            v-else-if="taskCurves[task.id]"
                            :points="task.points"
                            :curve="taskCurves[task.id]"
                            :targetPoint="{ x: task.target_x, y: task.correct_answer }"
                            :height="320"
                          />
                        </v-card>

                        <v-alert
                          v-if="taskExplanations[task.id]"
                          type="info"
                          variant="tonal"
                          color="secondary"
                          class="mt-4 rounded-lg"
                        >
                          <div class="text-body-2 leading-relaxed" v-html="taskExplanations[task.id]"></div>
                        </v-alert>
                      </div>
                    </v-expand-transition>
                    
                    <v-expand-transition>
                      <div v-if="results[task.id] === false" class="mt-4 text-error text-caption d-flex align-center">
                        <v-icon size="16" class="mr-1">mdi-alert-circle-outline</v-icon>
                        Попробуйте еще раз. Проверьте расчеты или метод.
                      </div>
                    </v-expand-transition>
                  </v-card>
                </div>
              </v-window-item>
            </v-window>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar.open" :color="snackbar.color" timeout="3500">
      {{ snackbar.message }}
    </v-snackbar>
  </v-container>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, defineComponent, h } from 'vue'
import axios from 'axios'
import tasksData from '../tasks.json'
import ChartView from '../components/ChartView.vue'
import renderMathInElement from 'katex/contrib/auto-render'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const tab = ref('theory')
const tasks = ref(tasksData)
const userAnswers = ref({})
const results = ref({})
const taskCurves = ref({})
const taskLoading = ref({})
const taskExplanations = ref({})
const theoryRoot = ref(null)
const snackbar = ref({ open: false, message: '', color: 'error' })

/**
 * Показывает короткое сообщение пользователю.
 * @param {string} message
 * @param {'error'|'warning'|'success'|'info'} [color='error']
 */
const notify = (message, color = 'error') => {
  snackbar.value = { open: true, message, color }
}

/**
 * Рендерит формулы KaTeX внутри теоретического блока.
 * Важно вызывать после изменения DOM (nextTick), иначе формулы могут не найтись.
 * @returns {Promise<void>}
 */
const renderTheoryMath = async () => {
  if (tab.value !== 'theory') return
  await nextTick()
  if (!theoryRoot.value) return
  renderMathInElement(theoryRoot.value, {
    delimiters: [
      { left: '$$', right: '$$', display: true },
      { left: '\\[', right: '\\]', display: true },
      { left: '\\(', right: '\\)', display: false },
      { left: '$', right: '$', display: false }
    ],
    throwOnError: false
  })
}

watch(tab, () => {
  renderTheoryMath()
})

/**
 * Сохраняет прогресс по заданиям в localStorage.
 * @returns {void}
 */
const saveResults = () => {
  const payload = { results: results.value, answers: userAnswers.value }
  localStorage.setItem('interpolation_practice_results', JSON.stringify(payload))
}

/**
 * Загружает интерполяционную кривую для задания через backend API.
 * @param {number} taskId
 * @returns {Promise<void>}
 */
const fetchTaskCurve = async (taskId) => {
  const task = tasks.value.find(t => t.id === taskId)
  if (!task) return

  taskLoading.value[taskId] = true
  try {
    const resp = await axios.post(`${API_URL}/interpolate`, {
      method: task.method,
      points: task.points,
      target_x: task.target_x
    })
    taskCurves.value[taskId] = Array.isArray(resp.data?.curve) ? resp.data.curve : []
  } catch (err) {
    taskCurves.value[taskId] = []
    notify('Не удалось загрузить кривую для графика. Попробуйте ещё раз.', 'warning')
  } finally {
    taskLoading.value[taskId] = false
  }
}

/**
 * Создаёт Vue-компонент иллюстрации на базе заранее подготовленного SVG.
 * @param {string} name
 * @param {string} svgMarkup
 * @returns {import('vue').Component}
 */
const createSvgIllustrationComponent = (name, svgMarkup) =>
  defineComponent({
    name,
    setup() {
      return () => h('div', { class: 'illustration-svg', innerHTML: svgMarkup })
    }
  })

const IllustrationPointsCurve = createSvgIllustrationComponent(
  'IllustrationPointsCurve',
  `
<svg viewBox="0 0 360 220" width="100%" height="220" role="img" aria-label="Точки и интерполяционная кривая">
  <rect x="0" y="0" width="360" height="220" rx="14" fill="#FFFFFF"/>
  <g stroke="#E2E8F0" stroke-width="1">
    <path d="M36 20V194M36 194H336" />
    <path d="M36 158H336" />
    <path d="M36 122H336" />
    <path d="M36 86H336" />
    <path d="M36 50H336" />
  </g>
  <path d="M56 156 C 110 40, 160 200, 214 96 S 304 68, 324 98" fill="none" stroke="#5E9DC8" stroke-width="4" stroke-linecap="round"/>
  <g fill="#E88CA5">
    <circle cx="56" cy="156" r="6"/>
    <circle cx="128" cy="74" r="6"/>
    <circle cx="214" cy="96" r="6"/>
    <circle cx="288" cy="78" r="6"/>
    <circle cx="324" cy="98" r="6"/>
  </g>
</svg>
  `.trim()
)

const IllustrationLinear = createSvgIllustrationComponent(
  'IllustrationLinear',
  `
<svg viewBox="0 0 360 220" width="100%" height="220" role="img" aria-label="Линейная интерполяция">
  <rect x="0" y="0" width="360" height="220" rx="14" fill="#FFFFFF"/>
  <g stroke="#E2E8F0" stroke-width="1">
    <path d="M36 20V194M36 194H336" />
  </g>
  <path d="M72 160 L 300 70" fill="none" stroke="#5E9DC8" stroke-width="4" stroke-linecap="round"/>
  <g fill="#E88CA5">
    <circle cx="72" cy="160" r="7"/>
    <circle cx="300" cy="70" r="7"/>
  </g>
  <g fill="#3BAB7B">
    <circle cx="186" cy="115" r="7"/>
  </g>
  <text x="186" y="142" text-anchor="middle" dominant-baseline="middle" font-family="Inter, sans-serif" font-size="12" fill="#64748B">промежуточная точка</text>
</svg>
  `.trim()
)

const IllustrationApplications = createSvgIllustrationComponent(
  'IllustrationApplications',
  `
<svg viewBox="0 0 360 220" width="100%" height="220" role="img" aria-label="Применение интерполяции">
  <rect x="0" y="0" width="360" height="220" rx="14" fill="#FFFFFF"/>
  <g fill="#F1F0F5" stroke="#E2E8F0" stroke-width="1">
    <rect x="24" y="30" width="312" height="160" rx="12"/>
  </g>
  <g stroke="#E2E8F0" stroke-width="1">
    <path d="M44 156H316" />
    <path d="M44 120H316" />
    <path d="M44 84H316" />
  </g>
  <g>
    <circle cx="78" cy="146" r="8" fill="#E88CA5"/>
    <circle cx="140" cy="110" r="8" fill="#E88CA5"/>
    <circle cx="210" cy="132" r="8" fill="#E88CA5"/>
    <circle cx="284" cy="92" r="8" fill="#E88CA5"/>
  </g>
  <path d="M78 146 C 110 120, 120 108, 140 110 S 188 150, 210 132 S 254 82, 284 92" fill="none" stroke="#5E9DC8" stroke-width="4" stroke-linecap="round"/>
  <g font-family="Inter, sans-serif" font-size="12" fill="#64748B">
    <text x="44" y="204" dominant-baseline="middle">датчики</text>
    <text x="140" y="204" dominant-baseline="middle">табличные данные</text>
    <text x="270" y="204" dominant-baseline="middle">графики</text>
  </g>
</svg>
  `.trim()
)

const IllustrationDividedDiff = createSvgIllustrationComponent(
  'IllustrationDividedDiff',
  `
<svg viewBox="0 0 360 220" width="100%" height="220" role="img" aria-label="Разделённые разности">
  <rect x="0" y="0" width="360" height="220" rx="14" fill="#FFFFFF"/>
  <g fill="#F1F0F5" stroke="#E2E8F0" stroke-width="1">
    <rect x="36" y="42" width="288" height="148" rx="10"/>
  </g>
  <g stroke="#E2E8F0" stroke-width="1">
    <path d="M36 78H324M36 114H324M36 150H324" />
    <path d="M108 42V190M180 42V190M252 42V190" />
  </g>
  <g font-family="JetBrains Mono, monospace" font-size="12" fill="#1E293B">
    <text x="54" y="66">x</text>
    <text x="126" y="66">f</text>
    <text x="198" y="66">Δ1</text>
    <text x="270" y="66">Δ2</text>
    <text x="54" y="102">0</text>
    <text x="126" y="102">1</text>
    <text x="54" y="138">1</text>
    <text x="126" y="138">3</text>
    <text x="54" y="174">2</text>
    <text x="126" y="174">2</text>
  </g>
  <path d="M196 116 L 268 152" stroke="#5E9DC8" stroke-width="3" stroke-linecap="round"/>
  <path d="M196 152 L 268 116" stroke="#E88CA5" stroke-width="3" stroke-linecap="round"/>
</svg>
  `.trim()
)

const IllustrationRunge = createSvgIllustrationComponent(
  'IllustrationRunge',
  `
<svg viewBox="0 0 360 220" width="100%" height="220" role="img" aria-label="Эффект Рунге">
  <rect x="0" y="0" width="360" height="220" rx="14" fill="#FFFFFF"/>
  <g stroke="#E2E8F0" stroke-width="1">
    <path d="M36 20V194M36 194H336" />
    <path d="M36 110H336" />
  </g>
  <path d="M48 130 C 90 90, 120 70, 160 90 S 240 150, 300 130" fill="none" stroke="#5E9DC8" stroke-width="3" stroke-linecap="round"/>
  <path d="M48 130 C 90 40, 130 190, 160 70 S 220 170, 240 60 S 300 200, 324 120" fill="none" stroke="#E88CA5" stroke-width="3" stroke-linecap="round"/>
  <g font-family="Inter, sans-serif" font-size="12" fill="#64748B">
    <path d="M52 206H92" stroke="#5E9DC8" stroke-width="3" stroke-linecap="round"/>
    <text x="100" y="206" dominant-baseline="middle">истинная функция</text>
    <path d="M214 206H254" stroke="#E88CA5" stroke-width="3" stroke-linecap="round"/>
    <text x="262" y="206" dominant-baseline="middle">полином высокой степени</text>
  </g>
</svg>
  `.trim()
)

const IllustrationCompare = createSvgIllustrationComponent(
  'IllustrationCompare',
  `
<svg viewBox="0 0 360 220" width="100%" height="220" role="img" aria-label="Сравнение Лагранжа и Ньютона">
  <rect x="0" y="0" width="360" height="220" rx="14" fill="#FFFFFF"/>
  <g fill="#F1F0F5" stroke="#E2E8F0" stroke-width="1">
    <rect x="28" y="38" width="144" height="144" rx="12"/>
    <rect x="188" y="38" width="144" height="144" rx="12"/>
  </g>
  <g font-family="Inter, sans-serif" font-size="12" fill="#1E293B" font-weight="700">
    <text x="100" y="62" text-anchor="middle" dominant-baseline="middle">Лагранж</text>
    <text x="260" y="62" text-anchor="middle" dominant-baseline="middle">Ньютон</text>
  </g>
  <g stroke="#E88CA5" stroke-width="3" stroke-linecap="round">
    <path d="M50 150 C 70 90, 110 90, 130 150" fill="none"/>
    <circle cx="66" cy="118" r="4" fill="#E88CA5"/>
    <circle cx="98" cy="108" r="4" fill="#E88CA5"/>
    <circle cx="128" cy="148" r="4" fill="#E88CA5"/>
  </g>
  <g stroke="#5E9DC8" stroke-width="3" stroke-linecap="round">
    <path d="M210 150 L 230 120 L 250 135 L 270 105 L 290 120" fill="none"/>
    <circle cx="210" cy="150" r="4" fill="#5E9DC8"/>
    <circle cx="230" cy="120" r="4" fill="#5E9DC8"/>
    <circle cx="250" cy="135" r="4" fill="#5E9DC8"/>
    <circle cx="270" cy="105" r="4" fill="#5E9DC8"/>
    <circle cx="290" cy="120" r="4" fill="#5E9DC8"/>
  </g>
  <g font-family="Inter, sans-serif" font-size="11" fill="#64748B">
    <text x="100" y="194" text-anchor="middle" dominant-baseline="middle">единый полином</text>
    <text x="260" y="194" text-anchor="middle" dominant-baseline="middle">наращивание по шагам</text>
  </g>
</svg>
  `.trim()
)

const theorySections = [
  {
    id: 'what-is',
    title: 'Что такое интерполяция?',
    icon: 'mdi-book-open-variant',
    accentColor: 'primary',
    layout: 'stack',
    paragraphsHtml: [
      'Интерполяция - это способ восстановить значения <strong>между</strong> уже известными точками. Представьте, что у вас есть несколько измерений с датчика, таблица из эксперимента или точки на графике, а между ними есть пробелы. Интерполяция помогает аккуратно заполнить эти пробелы так, чтобы новая функция проходила через все заданные узлы.',
      'Иначе говоря, мы не пытаемся угадать поведение функции "из воздуха", а строим модель на основе уже известных данных. Это особенно удобно, когда точная формула неизвестна, слишком сложна или отсутствует, но есть надёжные измерения.',
      'Основная идея: нужно построить функцию \\(f(x)\\), которая проходит точно через точки \\((x_i, y_i)\\).'
    ],
    notes: [
      {
        title: 'Интерполяция - это не угадывание',
        bodyHtml:
          'По сути, мы "восстанавливаем недостающее" между измерениями. Это помогает получать промежуточные значения для расчётов, визуализации и анализа.'
      },
      {
        title: 'Интерполяция vs экстраполяция',
        bodyHtml:
          'Интерполяция работает <strong>внутри</strong> диапазона известных точек. Экстраполяция - <strong>за его пределами</strong>, и она обычно заметно менее надёжна.'
      }
    ],
    illustrationComponent: IllustrationPointsCurve
  },
  {
    id: 'where-used',
    title: 'Где это применяется',
    icon: 'mdi-map-marker-path',
    accentColor: 'secondary',
    layout: 'stack',
    paragraphsHtml: [
      'Интерполяция используется почти везде, где данные получаются не непрерывно, а отдельными измерениями. В инженерии она встречается в калибровке датчиков, обработке сигналов, робототехнике, авиации, строительстве и моделировании физических процессов.',
      'Если известны значения температуры в нескольких моментах, интерполяция помогает оценить температуру в промежуточный момент. Если есть табличные данные давления, скорости или яркости - можно восстановить плавную зависимость между ними.'
    ],
    notes: [
      {
        title: 'Часто работает "за кадром"',
        bodyHtml:
          'Интерполяция используется в графике, навигации, обработке сенсорных данных и даже при построении карт погоды - вы просто не замечаете её напрямую.'
      },
      {
        title: 'Почему это удобно в инженерии',
        bodyHtml:
          'Когда формула поведения системы неизвестна или неудобна, но есть таблица измерений, интерполяция превращает дискретные точки в удобную модель для расчётов.'
      }
    ],
    illustrationComponent: IllustrationApplications
  },
  {
    id: 'linear',
    title: 'Линейная интерполяция',
    icon: 'mdi-chart-timeline-variant',
    accentColor: 'primary',
    paragraphsHtml: [
      'Линейная интерполяция - самый простой и интуитивный способ оценки промежуточного значения. Мы предполагаем, что между двумя известными точками график идёт по прямой.',
      'Метод работает быстро, легко программируется и особенно удобен, когда нужно получить ориентировочную оценку без сложных вычислений. Но если зависимость явно искривлена, точность может заметно упасть.'
    ],
    formulaLatex: String.raw`<div>\[
y = y_0 + (x - x_0)\cdot \frac{y_1 - y_0}{x_1 - x_0}
\]</div>`,
    notes: [
      {
        title: 'Когда метод особенно хорош',
        bodyHtml:
          'Когда точки близко друг к другу и зависимость почти прямая, линейная интерполяция даёт надёжный результат при минимальных вычислениях.'
      },
      {
        title: 'Где встречается чаще всего',
        bodyHtml:
          'Калибровка датчиков, системы управления, графические и физические расчёты - это один из самых "рабочих" методов в инженерной практике.'
      }
    ],
    illustrationComponent: IllustrationLinear
  },
  {
    id: 'lagrange',
    title: 'Полином Лагранжа',
    icon: 'mdi-function-variant',
    accentColor: 'secondary',
    paragraphsHtml: [
      'Метод Лагранжа строит <strong>единый многочлен</strong>, который проходит через все заданные точки сразу. Это уже не кусочная оценка, а полиномиальная модель, описывающая весь набор точек одной формулой.',
      'Метод математически красивый и универсальный, но при большом числе точек полином может начать сильно колебаться, особенно возле краёв интервала - это известно как <strong>эффект Рунге</strong>.'
    ],
    formulaLatex: String.raw`<div>\[
P(x)=\sum_{i=0}^{n} y_i \prod_{j=0, j\neq i}^{n}\frac{x-x_j}{x_i-x_j}
\]</div>`,
    notes: [
      {
        title: 'Про эффект Рунге',
        bodyHtml:
          'Чем выше степень полинома, тем сильнее он может "раскачиваться" на краях интервала, особенно на равномерной сетке узлов.'
      },
      {
        title: 'Немного истории',
        bodyHtml:
          'Метод назван в честь Жозефа Луи Лагранжа - одного из крупнейших математиков XVIII века. Его имя встречается во многих разделах математики и физики.'
      }
    ],
    illustrationComponent: IllustrationCompare
  },
  {
    id: 'newton',
    title: 'Полином Ньютона',
    icon: 'mdi-table-large',
    accentColor: 'primary',
    paragraphsHtml: [
      'Полином Ньютона тоже строит интерполяционную кривую через заданные точки, но делает это через <strong>разделённые разности</strong>. Его форма удобна тем, что каждый следующий член добавляется постепенно.',
      'Главное преимущество: при добавлении новой точки не нужно пересчитывать всё с самого начала - достаточно дописать новый член полинома.'
    ],
    formulaLatex: String.raw`<div>\[
P_n(x)=f(x_0)+(x-x_0)f[x_0,x_1]+(x-x_0)(x-x_1)f[x_0,x_1,x_2]+\dots
\]</div>`,
    notes: [
      {
        title: 'Почему это "инженерно" удобно',
        bodyHtml:
          'Если данные поступают постепенно, Ньютона удобно расширять по шагам - старая часть расчёта остаётся, а новая просто добавляется.'
      },
      {
        title: 'Немного истории',
        bodyHtml:
          'Подход основан на идеях Исаака Ньютона и таблицах разделённых разностей, поэтому его часто используют в табличных вычислениях.'
      }
    ],
    illustrationComponent: IllustrationDividedDiff
  },
  {
    id: 'choice',
    title: 'Что выбрать: Лагранж или Ньютон',
    icon: 'mdi-scale-balance',
    accentColor: 'secondary',
    paragraphsHtml: [
      'Оба метода дают один и тот же интерполяционный полином, если построены по одним и тем же точкам. Разница не в результате, а в форме записи и удобстве вычислений.',
      'Лагранж хорош как математически прозрачный способ записи, а Ньютон - как более гибкий вариант для поэтапного расширения набора данных.'
    ],
    notes: [
      {
        title: 'Важно помнить про границы данных',
        bodyHtml:
          'Если нужно выйти за пределы известных точек, это уже экстраполяция - и здесь ошибка обычно растёт быстрее.'
      }
    ],
    illustrationComponent: IllustrationCompare
  },
  {
    id: 'runge',
    title: 'Эффект Рунге',
    icon: 'mdi-waveform',
    accentColor: 'primary',
    paragraphsHtml: [
      'Эффект Рунге - пример того, что в интерполяции не всегда работает правило "чем больше точек, тем лучше". При использовании полиномов высокой степени на равномерной сетке кривая может начать сильно колебаться на краях интервала.',
      'На практике поэтому часто выбирают более устойчивые подходы: меньшее число точек, другой выбор узлов или альтернативные методы аппроксимации.'
    ],
    notes: [
      {
        title: 'Что это даёт на практике',
        bodyHtml:
          'Иногда более простой метод даёт более надёжный результат, чем "более умный" на бумаге - особенно при шумных данных и больших наборах узлов.'
      }
    ],
    illustrationComponent: IllustrationRunge
  }
]

const triggerConfetti = async () => {
  try {
    const mod = await import('canvas-confetti')
    const confetti = mod.default ?? mod
    confetti({
      particleCount: 140,
      spread: 70,
      startVelocity: 35,
      origin: { y: 0.7 }
    })
  } catch {}
}

/**
 * Проверяет ответ пользователя по заданию с учётом допустимой погрешности.
 * @param {number} id
 * @returns {void}
 */
const checkAnswer = (id) => {
  const task = tasks.value.find(t => t.id === id)
  if (!task) return

  const raw = userAnswers.value[id]
  const userVal = Number(raw)
  if (!Number.isFinite(userVal)) return

  const precision = Number.isFinite(Number(task.precision)) ? Number(task.precision) : 0.01
  const prev = results.value[id]
  const isCorrect = Math.abs(userVal - task.correct_answer) <= precision

  results.value[id] = isCorrect
  saveResults()

  if (isCorrect) {
    taskExplanations.value[id] = buildTaskExplanation(task, userVal)
    if (prev !== true) triggerConfetti()
    fetchTaskCurve(id)
  }
}

const getAnswerColor = (id) => {
  if (results.value[id] === undefined) return 'primary'
  return results.value[id] ? 'success' : 'error'
}

const getAnswerIcon = (id) => {
  if (results.value[id] === undefined) return ''
  return results.value[id] ? 'mdi-check-circle' : 'mdi-alert-circle'
}

/**
 * Формирует короткое объяснение решения, соответствующее выбранному методу.
 * @param {{method: string, target_x: number, correct_answer: number}} task
 * @param {number} userValue
 * @returns {string}
 */
const buildTaskExplanation = (task, userValue) => {
  const method = task.method
  const x = Number(task.target_x)
  const y = Number(task.correct_answer)
  const shownUser = Number.isFinite(userValue) ? userValue : y

  if (method === 'linear') {
    return `Ответ считается верным с учётом погрешности. Для линейной интерполяции берём две ближайшие точки и считаем значение по прямой. В точке <span class="math-font">x = ${x}</span> получаем <span class="math-font">f(x) ≈ ${y}</span>. Ваш ответ: <span class="math-font">${shownUser}</span>.`
  }

  const methodName = method === 'lagrange' ? 'Лагранжа' : 'Ньютона'
  return `Здесь применяется полиномиальная интерполяция методом ${methodName}. Сервис строит интерполяционный полином по заданным точкам и вычисляет значение в точке <span class="math-font">x = ${x}</span>: <span class="math-font">f(x) ≈ ${y}</span>. Ваш ответ: <span class="math-font">${shownUser}</span>.`
}

/**
 * Загружает сохранённый прогресс и, если нужно, догружает графики для выполненных заданий.
 * @returns {Promise<void>}
 */
const loadSavedProgress = async () => {
  const saved = localStorage.getItem('interpolation_practice_results')
  if (!saved) return

  try {
    const parsed = JSON.parse(saved)
    results.value = parsed?.results && typeof parsed.results === 'object' ? parsed.results : {}
    userAnswers.value = parsed?.answers && typeof parsed.answers === 'object' ? parsed.answers : {}
  } catch {
    localStorage.removeItem('interpolation_practice_results')
    results.value = {}
    userAnswers.value = {}
    notify('Сохранённый прогресс был повреждён и сброшен.', 'warning')
  }

  const completedIds = Object.keys(results.value)
    .map(v => Number(v))
    .filter(v => Number.isFinite(v) && results.value[String(v)] === true)

  await Promise.allSettled(completedIds.map(id => fetchTaskCurve(id)))
}

onMounted(async () => {
  await loadSavedProgress()
  await renderTheoryMath()
})
</script>

<style scoped>
.leading-relaxed {
  line-height: 1.625;
}
.bg-surface-light {
  background-color: #fcfcfd !important;
}
.theory-card {
  border: 1px solid #E2E8F0;
}
.formula-block {
  background: var(--color-rem-blue-light);
  border-color: rgba(94, 157, 200, 0.25) !important;
  padding: 12px 14px;
}
.illustration-card {
  background: #FFFFFF;
  border: 1px solid #E2E8F0;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: var(--shadow-card);
}
:deep(.illustration-svg svg) {
  display: block;
  width: 100%;
  height: auto;
}
.note-panel {
  border: 1px solid #E2E8F0;
  border-radius: 12px;
  overflow: hidden;
}
.notes-panels :deep(.v-expansion-panel)::before {
  opacity: 0 !important;
  border-top: 0 !important;
}
.notes-panels :deep(.v-expansion-panel-title) {
  border-top: 0 !important;
}
.notes-panels :deep(.v-expansion-panel-title)::before {
  opacity: 0 !important;
}
.notes-panels :deep(.v-expansion-panels--variant-accordion) {
  border-top: 0 !important;
}
.notes-panels :deep(.v-expansion-panel) {
  margin-top: 12px;
}
.notes-panels :deep(.v-expansion-panel:first-child) {
  margin-top: 0;
}
</style>
