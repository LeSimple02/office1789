<script setup>


//import {ref} from "vue"
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction'

let calendarOptions = {
        plugins: [ dayGridPlugin, interactionPlugin ],
        initialView: 'dayGridMonth',
        events: [
          { title: 'event 1', date: '2025-05-01' },
          { title: 'event 2', date: '2019-04-02' }
        ]
}


</script>
<template>
	<FullCalendar :options="calendarOptions" />
</template>


<style scoped>


</style>
