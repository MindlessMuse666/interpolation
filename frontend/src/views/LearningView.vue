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
              <!-- Theory -->
              <v-window-item value="theory">
                <div v-for="(t, index) in theory" :key="t.title" class="mb-10">
                  <div class="d-flex align-center mb-4">
                    <v-icon :color="index % 2 === 0 ? 'primary' : 'secondary'" size="32" class="mr-3">
                      {{ index % 2 === 0 ? 'mdi-book-open-variant' : 'mdi-lightbulb-on-outline' }}
                    </v-icon>
                    <h2 class="text-h4 font-weight-bold">{{ t.title }}</h2>
                  </div>
                  <p class="text-body-1 text-text-primary leading-relaxed mb-4">{{ t.content }}</p>
                  <v-divider v-if="index < theory.length - 1" class="mt-8"></v-divider>
                </div>
              </v-window-item>

              <!-- Practice -->
              <v-window-item value="practice">
                <v-alert
                  color="secondary"
                  variant="tonal"
                  class="mb-8 rounded-lg"
                  border="start"
                >
                  <template v-slot:prepend>
                    <v-icon size="24">mdi-information-outline</v-icon>
                  </template>
                  Решите задачи, используя изученные методы. Введите ответ с точностью до 0.01.
                </v-alert>

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
                    </v-row>
                    
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
  </v-container>
</template>

<script setup>
import { ref } from 'vue'
import tasksData from '../tasks.json'

const tab = ref('theory')
const tasks = ref(tasksData)
const userAnswers = ref({})
const results = ref({})

const theory = [
  {
    title: 'Что такое интерполяция?',
    content: 'Интерполяция — это математический метод нахождения промежуточных значений величины по имеющемуся дискретному набору известных данных. В инженерных расчетах это позволяет восстановить функцию по отдельным точкам измерений.'
  },
  {
    title: 'Линейная интерполяция',
    content: 'Самый простой и интуитивный метод. Мы предполагаем, что между двумя известными точками функция ведет себя как прямая линия. Формула для расчета: y = y0 + (x - x0) * (y1 - y0) / (x1 - x0).'
  },
  {
    title: 'Полином Лагранжа',
    content: 'Этот метод строит единый многочлен, который проходит ровно через все заданные точки. Он элегантен математически, но при большом количестве точек может давать сильные осцилляции на краях интервала (эффект Рунге).'
  },
  {
    title: 'Полином Ньютона',
    content: 'В отличие от Лагранжа, метод Ньютона строится итеративно с помощью разделенных разностей. Главное преимущество: при добавлении новой точки данных не нужно пересчитывать весь полином заново — достаточно добавить одно новое слагаемое.'
  }
]

const checkAnswer = (id) => {
  const task = tasks.value.find(t => t.id === id)
  const userVal = userAnswers.value[id]
  if (userVal === undefined || userVal === '') return

  const isCorrect = Math.abs(userVal - task.correct_answer) < 0.01
  results.value[id] = isCorrect
}

const getAnswerColor = (id) => {
  if (results.value[id] === undefined) return 'primary'
  return results.value[id] ? 'success' : 'error'
}

const getAnswerIcon = (id) => {
  if (results.value[id] === undefined) return ''
  return results.value[id] ? 'mdi-check-circle' : 'mdi-alert-circle'
}
</script>

<style scoped>
.leading-relaxed {
  line-height: 1.625;
}
.bg-surface-light {
  background-color: #fcfcfd !important;
}
</style>
