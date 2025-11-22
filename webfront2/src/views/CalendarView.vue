<script setup>
import { useI18n } from 'vue-i18n'
import { computed, ref, onMounted } from 'vue'
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction'
import { gls } from '../stores/global'

const { locale } = useI18n()

const showAddEventModal = ref(false)
const showDeleteEventModal = ref(false)
const selectedEvent = ref(null)
const newEventTitle = ref('')
const newEventDate = ref('')
const newEventStartTime = ref('09:00')
const newEventEndTime = ref('10:00')
const newEventDescription = ref('')
const newEventLocation = ref('')
const events = ref([])

// Charger les événements depuis le backend
async function loadEvents() {
  try {
    const response = await fetch('http://localhost:8080/api/calendar/events/get', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: gls().username,
        token: gls().sessionT
      })
    })
    if (response.ok) {
      const data = await response.json()
      // Transformer les événements pour FullCalendar
      events.value = (data.events || []).map(event => ({
        id: event.id,
        title: event.title,
        start: event.start,
        end: event.end,
        extendedProps: {
          description: event.description,
          location: event.location
        }
      }))
    }
  } catch (error) {
    console.error('Erreur lors du chargement des événements:', error)
  }
}

// Ajouter un événement
async function addEvent() {
  if (newEventTitle.value && newEventDate.value) {
    const newEvent = {
      username: gls().username,
      token: gls().sessionT,
      title: newEventTitle.value,
      start: `${newEventDate.value}T${newEventStartTime.value}:00Z`,
      end: `${newEventDate.value}T${newEventEndTime.value}:00Z`,
      description: newEventDescription.value,
      location: newEventLocation.value
    }
    
    console.log('Envoi événement:', newEvent)
    
    try {
      const response = await fetch('http://localhost:8080/api/calendar/events/create', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newEvent)
      })
      
      const result = await response.json()
      console.log('Réponse:', result)
      
      if (response.ok) {
        await loadEvents()
        resetForm()
        showAddEventModal.value = false
      } else {
        console.error('Erreur serveur:', result)
      }
    } catch (error) {
      console.error('Erreur lors de l\'ajout de l\'événement:', error)
    }
  }
}

// Supprimer un événement
async function deleteEvent() {
  if (selectedEvent.value) {
    try {
      const response = await fetch('http://localhost:8080/api/calendar/events/delete', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username: gls().username,
          token: gls().sessionT,
          event_id: selectedEvent.value.id.toString()
        })
      })
      
      if (response.ok) {
        await loadEvents()
        showDeleteEventModal.value = false
        selectedEvent.value = null
      }
    } catch (error) {
      console.error('Erreur lors de la suppression de l\'événement:', error)
    }
  }
}

// Réinitialiser le formulaire
function resetForm() {
  newEventTitle.value = ''
  newEventDate.value = ''
  newEventStartTime.value = '09:00'
  newEventEndTime.value = '10:00'
  newEventDescription.value = ''
  newEventLocation.value = ''
}

// Gérer le clic sur un événement
function handleEventClick(info) {
  selectedEvent.value = info.event
  showDeleteEventModal.value = true
}

const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, interactionPlugin],
  initialView: 'dayGridMonth',
  locale: locale.value,
  events: events.value,
  eventClick: handleEventClick
}))

onMounted(() => {
  loadEvents()
})
</script>

<template>
  <div class="calendar-view">
    <!-- Hero Card -->
    <div class="hero-card">
      <div class="calendar-icon-wrapper">
        <svg class="calendar-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
        </svg>
      </div>
      <h1>📅 {{ $t('calendarOffice1789') }}</h1>
      <h2>{{ $t('organizeAppointments') }}</h2>
      
      <button @click="showAddEventModal = true" class="btn-add-event">
        <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M12 4v16m8-8H4"/>
        </svg>
        {{ $t('addEvent') }}
      </button>
    </div>

    <!-- Calendar -->
    <div class="calendar-wrapper">
      <FullCalendar :options="calendarOptions" />
    </div>

    <!-- Modal Add Event -->
    <div v-if="showAddEventModal" class="modal-overlay" @click.self="showAddEventModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>✨ {{ $t('newEvent') }}</h3>
          <button class="close-btn" @click="showAddEventModal = false">✕</button>
        </div>
        
        <div class="form-grid">
          <div class="form-group">
            <label>{{ $t('title') }} *</label>
            <input v-model="newEventTitle" class="input" :placeholder="$t('eventTitle')" />
          </div>
          
          <div class="form-group">
            <label>{{ $t('date') }} *</label>
            <input v-model="newEventDate" type="date" class="input" />
          </div>
          
          <div class="form-group">
            <label>{{ $t('startTime') }}</label>
            <input v-model="newEventStartTime" type="time" class="input" />
          </div>
          
          <div class="form-group">
            <label>{{ $t('endTime') }}</label>
            <input v-model="newEventEndTime" type="time" class="input" />
          </div>
          
          <div class="form-group full-width">
            <label>{{ $t('location') }}</label>
            <input v-model="newEventLocation" class="input" :placeholder="$t('eventLocation')" />
          </div>
          
          <div class="form-group full-width">
            <label>{{ $t('description') }}</label>
            <textarea v-model="newEventDescription" class="textarea" :placeholder="$t('eventDescription')" rows="4"></textarea>
          </div>
        </div>
        
        <div class="modal-actions">
          <button class="btn btn-primary" @click="addEvent">
            <svg width="20" height="20" fill="currentColor" viewBox="0 0 24 24">
              <path d="M12 4v16m8-8H4"/>
            </svg>
            {{ $t('add') }}
          </button>
          <button class="btn btn-secondary" @click="showAddEventModal = false">{{ $t('cancel') }}</button>
        </div>
      </div>
    </div>

    <!-- Modal Delete Event -->
    <div v-if="showDeleteEventModal" class="modal-overlay" @click.self="showDeleteEventModal = false">
      <div class="modal modal-small">
        <div class="modal-header">
          <h3>🗑️ {{ $t('deleteEvent') }}</h3>
          <button class="close-btn" @click="showDeleteEventModal = false">✕</button>
        </div>
        
        <div class="modal-body">
          <p class="delete-message">
            {{ $t('confirmDeleteEvent') }}
            <strong>"{{ selectedEvent?.title }}"</strong> ?
          </p>
        </div>
        
        <div class="modal-actions">
          <button class="btn btn-danger" @click="deleteEvent">
            <svg width="20" height="20" fill="currentColor" viewBox="0 0 24 24">
              <path d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/>
            </svg>
            {{ $t('delete') }}
          </button>
          <button class="btn btn-secondary" @click="showDeleteEventModal = false">{{ $t('cancel') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.calendar-view {
  width: 100%;
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Hero Card */
.hero-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  width: 100%;
  padding: 48px 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
  gap: 24px;
}

.dark .hero-card {
  background: #1C1C1E;
}

.calendar-icon-wrapper {
  width: 100px;
  height: 100px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.calendar-icon {
  width: 50px;
  height: 50px;
  color: white;
}

.hero-card h1 {
  font-family: 'Roboto', sans-serif;
  font-size: 3rem;
  font-weight: 700;
  text-align: center;
  letter-spacing: 2px;
  color: #222;
  margin: 0;
}

.dark .hero-card h1 {
  color: white;
}

.hero-card h2 {
  font-family: 'Roboto', sans-serif;
  font-size: 1.5rem;
  font-style: italic;
  text-align: center;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.5px;
  margin: 0;
}

.btn-add-event {
  font-family: 'Roboto', sans-serif;
  font-size: 1.2rem;
  padding: 1rem 2.5rem;
  border-radius: 32px;
  border: none;
  cursor: pointer;
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  font-weight: 600;
}

.btn-add-event:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.btn-icon {
  width: 24px;
  height: 24px;
}

.calendar-wrapper {
  margin-bottom: 32px;
}

.dark .calendar-wrapper :deep(*:not(.fc-toolbar):not(.fc-toolbar *):not(.fc-button):not(.fc-button *):not(.fc-toolbar-title)) {
  background-color: #000000 !important;
  background: #000000 !important;
}

/* Style moderne pour FullCalendar */
:deep(.fc) {
  font-family: 'Roboto', Arial, sans-serif;
  background: rgba(245,245,247,0.85);
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.12);
  padding: 24px;
}

.dark :deep(.fc),
.dark :deep(.fc *),
.dark :deep(.fc-view),
.dark :deep(.fc-view-harness),
.dark :deep(.fc-daygrid),
.dark :deep(.fc-daygrid-body),
.dark :deep(.fc-scrollgrid),
.dark :deep(.fc-scrollgrid-sync-table),
.dark :deep(.fc-daygrid-body-unbalanced),
.dark :deep(.fc-daygrid-body-natural),
.dark :deep(.fc-scroller),
.dark :deep(.fc-scroller-liquid-absolute),
.dark :deep(td),
.dark :deep(th),
.dark :deep(.fc-scrollgrid-section),
.dark :deep(.fc-scrollgrid-section-body),
.dark :deep(.fc-scrollgrid-section-header) {
  background-color: #000000 !important;
  background: #000000 !important;
  color: white !important;
}

.dark :deep(.fc-daygrid-day-number),
.dark :deep(.fc-col-header-cell-cushion),
.dark :deep(.fc-daygrid-day-top) {
  color: white !important;
}

.dark :deep(.fc-daygrid-day),
.dark :deep(.fc-daygrid-day-frame),
.dark :deep(.fc-daygrid-day-bg) {
  background: #000000 !important;
  background-color: #000000 !important;
  border-color: rgba(255, 255, 255, 0.1) !important;
}

.dark :deep(.fc-theme-standard td),
.dark :deep(.fc-theme-standard th),
.dark :deep(.fc-col-header-cell) {
  border-color: rgba(255, 255, 255, 0.1) !important;
  background: #000000 !important;
  background-color: #000000 !important;
}

:deep(.fc-toolbar) {
  background: linear-gradient(30deg, blue, red);
  font-family: roboto;
  border-radius: 16px;
  color: #fff;
  margin-bottom: 16px;
  padding: 12px 16px;
}

:deep(.fc-toolbar-title) {
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: 2px;
}

:deep(.fc-button) {
  background: #06c;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  padding: 6px 16px;
  margin: 0 4px;
  transition: background 0.2s, box-shadow 0.2s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.10);
  cursor: pointer;
}

:deep(.fc-button:hover),
:deep(.fc-button-active) {
  background: linear-gradient(30deg, blue, red);
  color: #fff;
  box-shadow: 0 4px 16px rgba(0,0,0,0.18);
}

:deep(.fc-daygrid-day) {
  border-radius: 8px;
  transition: background 0.2s;
}

:deep(.fc-daygrid-day:hover) {
  background: rgba(0,48,143,0.08);
}

:deep(.fc-event) {
  background: linear-gradient(30deg, blue, red);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.10);
  padding: 2px 8px;
}

:deep(.fc-daygrid-event-dot) {
  background: #06c;
}

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.modal {
  background: white;
  border-radius: 24px;
  padding: 0;
  width: 90%;
  max-width: 600px;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.3);
  animation: slideIn 0.3s ease;
  overflow: hidden;
}

.modal-small {
  max-width: 400px;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dark .modal {
  background: #1C1C1E;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px;
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
}

.modal-header h3 {
  font-family: 'Roboto', sans-serif;
  font-size: 1.5rem;
  margin: 0;
  font-weight: 700;
}

.close-btn {
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  font-size: 1.5rem;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: rotate(90deg);
}

.modal-body {
  padding: 32px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  padding: 32px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group.full-width {
  grid-column: 1 / -1;
}

.form-group label {
  font-family: 'Roboto', sans-serif;
  font-size: 0.95rem;
  font-weight: 600;
  color: #333;
}

.dark .form-group label {
  color: #eee;
}

.input {
  width: 100%;
  padding: 12px 16px;
  border-radius: 12px;
  border: 2px solid rgba(0, 48, 143, 0.2);
  font-size: 1rem;
  font-family: 'Roboto', sans-serif;
  background: white;
  box-sizing: border-box;
  transition: all 0.3s ease;
}

.input:focus {
  outline: none;
  border-color: blue;
  box-shadow: 0 0 0 3px rgba(0, 48, 143, 0.1);
}

.textarea {
  width: 100%;
  padding: 12px 16px;
  border-radius: 12px;
  border: 2px solid rgba(0, 48, 143, 0.2);
  font-size: 1rem;
  font-family: 'Roboto', sans-serif;
  background: white;
  box-sizing: border-box;
  transition: all 0.3s ease;
  resize: vertical;
  min-height: 100px;
}

.textarea:focus {
  outline: none;
  border-color: blue;
  box-shadow: 0 0 0 3px rgba(0, 48, 143, 0.1);
}

.dark .input,
.dark .textarea {
  background: rgba(30, 30, 40, 0.95);
  color: white;
  border-color: rgba(255, 255, 255, 0.2);
}

.dark .input:focus,
.dark .textarea:focus {
  border-color: rgba(255, 255, 255, 0.5);
  box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.1);
}

.delete-message {
  font-family: 'Roboto', sans-serif;
  font-size: 1.1rem;
  color: #333;
  text-align: center;
  line-height: 1.6;
}

.dark .delete-message {
  color: #eee;
}

.delete-message strong {
  color: #d32f2f;
  font-weight: 700;
}

.dark .delete-message strong {
  color: #ff6b6b;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  padding: 24px 32px;
  background: rgba(245, 245, 247, 0.5);
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}

.dark .modal-actions {
  background: rgba(30, 30, 40, 0.5);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.btn {
  padding: 12px 32px;
  border-radius: 16px;
  border: none;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: 'Roboto', sans-serif;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.btn svg {
  width: 20px;
  height: 20px;
}

.btn-primary {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.btn-danger {
  background: -webkit-linear-gradient(30deg, #d32f2f, #c62828);
  color: white;
}

.btn-danger:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(211, 47, 47, 0.4);
}

.btn-secondary {
  background: #eee;
  color: #333;
}

.btn-secondary:hover {
  background: #ddd;
}

.dark .btn-secondary {
  background: #333;
  color: #eee;
}

.dark .btn-secondary:hover {
  background: #444;
}

/* Responsive */
@media (max-width: 767px) {
  .hero-card {
    padding: 32px 20px;
  }

  .hero-card h1 {
    font-size: 1.8rem;
  }

  .hero-card h2 {
    font-size: 1rem;
  }

  .btn-add-event {
    width: 100%;
    justify-content: center;
    font-size: 1rem;
    padding: 0.875rem 1.5rem;
  }

  :deep(.fc) {
    padding: 8px;
    border-radius: 12px;
  }

  :deep(.fc-toolbar-title) {
    font-size: 1.1rem;
  }
}
</style>