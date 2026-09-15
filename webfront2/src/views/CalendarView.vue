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

async function loadEvents() {
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/calendar/events/get`, {
      method: 'POST', credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: gls().username, token: gls().sessionT })
    })
    if (response.ok) {
      const data = await response.json()
      events.value = (data.events || []).map(event => ({
        id: event.id, title: event.title, start: event.start, end: event.end,
        extendedProps: { description: event.description, location: event.location }
      }))
    }
  } catch (error) { console.error('Erreur lors du chargement des événements:', error) }
}

async function addEvent() {
  if (newEventTitle.value && newEventDate.value) {
    const newEvent = {
      username: gls().username, token: gls().sessionT,
      title: newEventTitle.value,
      start: `${newEventDate.value}T${newEventStartTime.value}:00Z`,
      end: `${newEventDate.value}T${newEventEndTime.value}:00Z`,
      description: newEventDescription.value, location: newEventLocation.value
    }
    try {
      const response = await fetch(`${import.meta.env.VITE_APP_API}/api/calendar/events/create`, {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newEvent)
      })
      if (response.ok) { await loadEvents(); resetForm(); showAddEventModal.value = false }
    } catch (error) { console.error('Erreur lors de l\'ajout:', error) }
  }
}

async function deleteEvent() {
  if (selectedEvent.value) {
    try {
      const response = await fetch(`${import.meta.env.VITE_APP_API}/api/calendar/events/delete`, {
        method: 'POST', credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: gls().username, token: gls().sessionT, event_id: selectedEvent.value.id.toString() })
      })
      if (response.ok) { await loadEvents(); showDeleteEventModal.value = false; selectedEvent.value = null }
    } catch (error) { console.error('Erreur lors de la suppression:', error) }
  }
}

function resetForm() {
  newEventTitle.value = ''; newEventDate.value = ''
  newEventStartTime.value = '09:00'; newEventEndTime.value = '10:00'
  newEventDescription.value = ''; newEventLocation.value = ''
}

function handleEventClick(info) { selectedEvent.value = info.event; showDeleteEventModal.value = true }

const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, interactionPlugin],
  initialView: 'dayGridMonth',
  locale: locale.value,
  events: events.value,
  eventClick: handleEventClick
}))

onMounted(() => { loadEvents() })
</script>

<template>
  <div class="cal-page">
    <div class="cal-head">
      <div>
        <h1>Calendrier</h1>
        <p>Gérez vos événements et rendez-vous.</p>
      </div>
      <button @click="showAddEventModal = true" class="btn-add">+ Événement</button>
    </div>

    <div class="cal-wrap">
      <FullCalendar :options="calendarOptions" />
    </div>

    <!-- Add modal -->
    <div v-if="showAddEventModal" class="overlay" @click.self="showAddEventModal = false">
      <div class="dialog">
        <div class="dialog-head">
          <h3>Nouvel événement</h3>
          <button @click="showAddEventModal = false" class="close">×</button>
        </div>
        <div class="dialog-body">
          <div class="grid2">
            <div class="fld">
              <label>Titre</label>
              <input v-model="newEventTitle" type="text" placeholder="Titre de l'événement" />
            </div>
            <div class="fld">
              <label>Date</label>
              <input v-model="newEventDate" type="date" />
            </div>
            <div class="fld">
              <label>Heure de début</label>
              <input v-model="newEventStartTime" type="time" />
            </div>
            <div class="fld">
              <label>Heure de fin</label>
              <input v-model="newEventEndTime" type="time" />
            </div>
            <div class="fld full">
              <label>Lieu</label>
              <input v-model="newEventLocation" type="text" placeholder="Lieu (optionnel)" />
            </div>
            <div class="fld full">
              <label>Description</label>
              <textarea v-model="newEventDescription" placeholder="Description (optionnelle)"></textarea>
            </div>
          </div>
        </div>
        <div class="dialog-foot">
          <button @click="showAddEventModal = false" class="btn-cancel">Annuler</button>
          <button @click="addEvent" class="btn-confirm">Ajouter</button>
        </div>
      </div>
    </div>

    <!-- Delete modal -->
    <div v-if="showDeleteEventModal" class="overlay" @click.self="showDeleteEventModal = false">
      <div class="dialog dialog-sm">
        <div class="dialog-head">
          <h3>Supprimer ?</h3>
          <button @click="showDeleteEventModal = false" class="close">×</button>
        </div>
        <div class="dialog-body">
          <p class="del-msg">Voulez-vous supprimer "{{ selectedEvent?.title }}" ?</p>
        </div>
        <div class="dialog-foot">
          <button @click="showDeleteEventModal = false" class="btn-cancel">Annuler</button>
          <button @click="deleteEvent" class="btn-confirm">Supprimer</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cal-page { max-width: var(--content-max); margin: 0 auto; padding: 2.5rem 1.5rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.cal-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1.5rem; gap: 1rem; flex-wrap: wrap; animation: fadeInDown .4s var(--ease) both; }
.cal-head h1 { font-size: 1.75rem; font-weight: 800; letter-spacing: -.025em; }
.cal-head p { font-size: .9375rem; color: var(--text-secondary); }
.btn-add { padding: .625rem 1.25rem; font-size: .875rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); white-space: nowrap; }
.btn-add:hover { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-sm); transform: translateY(-2px); }
.cal-wrap { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.25rem; animation: fadeInUp .4s var(--ease) .1s both; }

/* FullCalendar overrides */
.cal-wrap :deep(.fc) { font-family: var(--font-sans); }
.cal-wrap :deep(.fc-toolbar-title) { font-size: 1.125rem; font-weight: 700; }
.cal-wrap :deep(.fc-button) { background: var(--bg-primary) !important; color: var(--text-primary) !important; border: 1px solid var(--border-strong) !important; border-radius: var(--r-sm) !important; text-transform: capitalize !important; font-weight: 500 !important; font-size: .8125rem !important; padding: .375rem .75rem !important; transition: all var(--tb) !important; }
.cal-wrap :deep(.fc-button:hover) { border-color: var(--ink) !important; background: var(--bg-tertiary) !important; }
.cal-wrap :deep(.fc-button-active) { background: var(--ink) !important; color: var(--paper) !important; border-color: var(--ink) !important; }
.cal-wrap :deep(.fc-daygrid-day) { transition: background var(--tb); }
.cal-wrap :deep(.fc-daygrid-day:hover) { background: var(--bg-secondary); }
.cal-wrap :deep(.fc-col-header-cell-cushion) { color: var(--text-secondary); font-weight: 600; font-size: .75rem; text-transform: uppercase; letter-spacing: .03em; }
.cal-wrap :deep(.fc-daygrid-day-number) { color: var(--text-primary); font-weight: 500; }
.cal-wrap :deep(.fc-event) { background: var(--ink) !important; border: none !important; border-radius: var(--r-sm) !important; padding: 2px 6px !important; font-size: .75rem !important; font-weight: 500 !important; cursor: pointer; transition: opacity var(--tb) !important; }
.cal-wrap :deep(.fc-event:hover) { opacity: .8; }
.cal-wrap :deep(.fc-today) { background: var(--bg-tertiary) !important; }
.cal-wrap :deep(.fc .fc-scrollgrid) { border: none !important; }
.cal-wrap :deep(.fc th, .fc td) { border-color: var(--border) !important; }

/* Modal */
.overlay { position: fixed; inset: 0; z-index: 999; background: var(--bg-overlay); backdrop-filter: blur(6px); display: flex; align-items: center; justify-content: center; animation: fadeIn .2s var(--ease) both; }
.dialog { width: calc(100% - 32px); max-width: 520px; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); box-shadow: var(--sh-xl); animation: popIn .25s var(--ease-spring) both; max-height: 90vh; display: flex; flex-direction: column; }
.dialog-sm { max-width: 380px; }
.dialog-head { display: flex; align-items: center; justify-content: space-between; padding: 1.125rem 1.5rem; border-bottom: 1px solid var(--border); }
.dialog-head h3 { font-size: 1.0625rem; font-weight: 700; }
.close { width: 32px; height: 32px; display: flex; align-items: center; justify-content: center; font-size: 20px; line-height: 1; background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r-sm); cursor: pointer; color: var(--text-secondary); transition: all var(--tb); }
.close:hover { border-color: var(--ink); background: var(--bg-tertiary); transform: rotate(90deg); color: var(--text-primary); }
.dialog-body { padding: 1.5rem; overflow-y: auto; }
.grid2 { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
.fld { display: flex; flex-direction: column; }
.fld.full { grid-column: 1 / -1; }
.fld label { font-size: .8125rem; font-weight: 600; color: var(--text-secondary); margin-bottom: .375rem; }
.fld input, .fld textarea { width: 100%; padding: .625rem .875rem; font-size: .9375rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: border-color var(--tb), box-shadow var(--tb); font-family: inherit; }
.fld input:focus, .fld textarea:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.fld textarea { resize: vertical; min-height: 60px; }
.del-msg { font-size: .9375rem; color: var(--text-primary); }
.dialog-foot { display: flex; gap: .625rem; padding: 1.125rem 1.5rem; border-top: 1px solid var(--border); background: var(--bg-secondary); }
.btn-cancel { flex: 1; padding: .625rem; font-size: .875rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-cancel:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.btn-confirm { flex: 1; padding: .625rem; font-size: .875rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-confirm:hover { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-sm); }

@media (max-width: 640px) { .grid2 { grid-template-columns: 1fr; } }
</style>
