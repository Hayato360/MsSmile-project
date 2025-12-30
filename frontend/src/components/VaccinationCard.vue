<script setup>
import { ref, watch } from 'vue'
import { Edit } from 'lucide-vue-next'

const props = defineProps({
  vaccineType: {
    type: Object,
    required: true,
  },
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  readonly: {
    type: Boolean,
    default: false,
  },
  canEdit: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'save', 'edit'])

const form = ref({ ...props.modelValue })

watch(
  () => props.modelValue,
  (newVal) => {
    form.value = { ...newVal }
  },
  { deep: true }
)

const updateField = () => {
  emit('update:modelValue', form.value)
}

// History Logic
const setHistoryStatus = (status) => {
  if (props.readonly) return
  
  if (status === 'unknown') {
    form.value.IsHistoryUnknown = true
    form.value.IsPreviouslyVaccinated = false
    form.value.PreviousDoses = 0
    form.value.LastPreviousDateYear = null
  } else if (status === 'previously') {
    form.value.IsHistoryUnknown = false
    form.value.IsPreviouslyVaccinated = true
  } else {
    // never
    form.value.IsHistoryUnknown = false
    form.value.IsPreviouslyVaccinated = false
    form.value.PreviousDoses = 0
    form.value.LastPreviousDateYear = null
  }
  updateField()
}

// Current Pregnancy Logic
const setCurrentStatus = (isVaccinating) => {
  if (props.readonly) return
  
  // If switching to not vaccinating, clear dates
  if (!isVaccinating) {
    form.value.Dose1DateDuringPreg = null
    form.value.Dose2DateDuringPreg = null
    form.value.Dose3DateDuringPreg = null
  } else {
    // If switching to vaccinating, clear reason
    form.value.ReasonForNotVaccinating = ''
  }
  updateField()
}

// Check if currently vaccinating (has at least one date or intended)
// Actually, backend doesn't store "IsVaccinatingDuringPregnancy" bool, 
// it infers from dates or Reason.
// So we need a local UI state or infer logic.
// Let's infer: If Reason is empty, assumes intent to vaccinate or already vaccinated?
// Or simpler: Just showing dates implies "Yes", showing Reason implies "No".
// But user wants a Radio.
const isVaccinatingCurrent = ref(!form.value.ReasonForNotVaccinating)

watch(isVaccinatingCurrent, (newVal) => {
  setCurrentStatus(newVal)
})

</script>

<template>
  <div class="vaccine-card">
    <div class="card-header">
      <h4>{{ vaccineType.Name }}</h4>
      <button v-if="canEdit" type="button" @click="$emit('edit')" class="btn-icon-edit" title="แก้ไข">
        <Edit size="16" />
      </button>
    </div>

    <!-- History Section -->
    <div class="section history-section">
      <div class="radio-item">
        <label class="radio-label">
          <input 
            type="radio" 
            :name="`history-${vaccineType.ID}`" 
            :checked="form.IsPreviouslyVaccinated"
            @change="setHistoryStatus('previously')"
            :disabled="readonly"
          >
          เคยฉีด
        </label>
        <div class="sub-inputs" v-if="form.IsPreviouslyVaccinated">
          <div class="inline-input">
            <span>จำนวน</span>
            <input 
              type="number" 
              v-model.number="form.PreviousDoses" 
              class="input-sm"
              @input="updateField"
              :disabled="readonly"
            >
            <span>เข็ม</span>
          </div>
          <div class="inline-input">
            <span>ครั้งสุดท้าย วันที่</span>
            <input 
              type="date" 
              v-model="form.LastPreviousDateYear" 
              class="input-md"
              @input="updateField"
              :disabled="readonly"
            >
          </div>
        </div>
      </div>

      <div class="radio-item">
        <label class="radio-label">
          <input 
            type="radio" 
            :name="`history-${vaccineType.ID}`" 
            :checked="!form.IsPreviouslyVaccinated && !form.IsHistoryUnknown"
            @change="setHistoryStatus('never')"
            :disabled="readonly"
          >
          ไม่เคยฉีด
        </label>
      </div>

      <div class="radio-item">
        <label class="radio-label">
          <input 
            type="radio" 
            :name="`history-${vaccineType.ID}`" 
            :checked="form.IsHistoryUnknown"
            @change="setHistoryStatus('unknown')"
            :disabled="readonly"
          >
          ไม่ทราบ/ไม่แน่ใจ
        </label>
      </div>
    </div>

    <div class="divider"></div>

    <!-- Current Pregnancy Section -->
    <div class="section current-section">
      <div class="section-title">ในระหว่างการตั้งครรภ์นี้</div>
      
      <div class="radio-item">
        <label class="radio-label">
          <input 
            type="radio" 
            :name="`current-${vaccineType.ID}`" 
            :value="true"
            v-model="isVaccinatingCurrent"
            :disabled="readonly"
          >
          ฉีดวัคซีน
        </label>
        <div class="sub-inputs" v-if="isVaccinatingCurrent">
          <div class="inline-input">
            <span>ครั้งที่ 1 วันที่</span>
            <input 
              type="date" 
              v-model="form.Dose1DateDuringPreg" 
              class="input-md"
              @input="updateField"
              :disabled="readonly"
            >
          </div>
          <div class="inline-input">
            <span>ครั้งที่ 2 วันที่</span>
            <input 
              type="date" 
              v-model="form.Dose2DateDuringPreg" 
              class="input-md"
              @input="updateField"
              :disabled="readonly"
            >
          </div>
           <div class="inline-input">
            <span>ครั้งที่ 3 วันที่</span>
            <input 
              type="date" 
              v-model="form.Dose3DateDuringPreg" 
              class="input-md"
              @input="updateField"
              :disabled="readonly"
            >
          </div>
        </div>
      </div>

      <div class="radio-item">
        <label class="radio-label">
          <input 
            type="radio" 
            :name="`current-${vaccineType.ID}`" 
            :value="false"
            v-model="isVaccinatingCurrent"
            :disabled="readonly"
          >
          ไม่ฉีดวัคซีนในครรภ์นี้
        </label>
         <div class="sub-inputs" v-if="!isVaccinatingCurrent">
            <input 
              type="text" 
              v-model="form.ReasonForNotVaccinating" 
              placeholder="ระบุเหตุผล (เช่น แพ้, ปฏิเสธ)"
              class="input-full"
              @input="updateField"
              :disabled="readonly"
            >
         </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vaccine-card {
  background: #e5e7eb; /* Gray background as per image */
  border-radius: 1rem;
  overflow: hidden;
  border: 1px solid #d1d5db;
  display: flex;
  flex-direction: column;
}

.card-header {
  background: #10b981;
  padding: 0.75rem;
  text-align: center;
  font-weight: bold;
  color: white;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.btn-icon-edit {
    background: rgba(255, 255, 255, 0.2);
    border: 1px solid rgba(255, 255, 255, 0.4);
    border-radius: 4px;
    padding: 4px 8px;
    cursor: pointer;
    display: flex;
    align-items: center;
    color: white;
}
.btn-icon-edit:hover {
    background: rgba(255, 255, 255, 0.3);
    border-color: rgba(255, 255, 255, 0.6);
}

.card-header h4 {
    margin: 0;
    font-size: 1rem;
}

.section {
  padding: 1rem;
  background: white; /* Inner white box concept from image logic, or direct on gray? 
                      Image looks like white paper on gray board? 
                      Let's stick to simple clean layout first. 
                      Actually, image shows white "Inset" boxes. */
  margin: 0.5rem;
  border-radius: 0.5rem;
}

.history-section {
    min-height: 140px;
}

.current-section {
    background-color: #fef9c3; /* Light yellow for current pregnancy, matching image hint? Or just white. let's stick to white or light yellow if meaningful. Image bottom box is yellowish? Hard to tell. Let's start with white. */
    background: white;
}

.divider {
    height: 2px;
    background: #4b5563;
    margin: 0 1rem;
    position: relative;
}
.divider::before, .divider::after {
    content: '';
    position: absolute;
    top: -3px;
    width: 8px;
    height: 8px;
    background: white;
    border: 2px solid #4b5563;
    
}
.divider::before { left: 0; }
.divider::after { right: 0; }
/* Just kidding, the image shows a line with squares at ends. simple hr is fine for now. */


.radio-item {
  margin-bottom: 0.5rem;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  font-size: 0.95rem;
  font-weight: 500;
  color: #374151;
}

.radio-label input[type="radio"] {
    width: 1.2rem;
    height: 1.2rem;
    accent-color: #10b981;
}

.sub-inputs {
  margin-left: 1.8rem;
  margin-top: 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.inline-input {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
}

.input-sm {
  width: 50px;
  padding: 0.25rem;
  border: 1px solid #d1d5db;
  border-radius: 0.25rem;
  text-align: center;
}

.input-md {
  width: 130px;
  padding: 0.25rem;
  border: 1px solid #d1d5db;
  border-radius: 0.25rem;
}

.input-full {
  width: 100%;
  padding: 0.4rem;
  border: 1px solid #d1d5db;
  border-radius: 0.25rem;
}

.section-title {
    font-weight: bold;
    margin-bottom: 0.5rem;
    color: #10b981;
    text-decoration: underline;
}
</style>
