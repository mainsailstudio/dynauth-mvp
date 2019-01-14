<template>
    <label for="email">Email:</label>

    <input type="email" name="email" v-model.trim="email" placeholder="test@dynauth.io">
    <button class="next" v-on:click="getLocksFromEmail" title="Next!">
    <i class="fas fa-angle-double-right"></i></button>
    <label v-if="locksExist">Locks:&nbsp;&nbsp;<span class="locks">{{ lockString }}</span>&nbsp;&nbsp;&nbsp;<i class="fas fa-sync fa-xs" id="refresh-icon" v-on:click="generateLocksAndKeys"></i></label>
    <input v-if="locksExist" type="text" name="user-input" placeholder="loading..." id="user-input">
    <button v-if="locksExist" class="next" v-on:click="authenticate" title="Next!"><i class="fas fa-angle-double-right"></i></button>
</template>

<script>
export default {
  data () {
    return {}
  }
}
</script>
