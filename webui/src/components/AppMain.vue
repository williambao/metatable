<template>
    <section class="app-main">
        <transition mode="out-in" enter-active-class="fadeIn" leave-active-class="fadeOut" appear>
            <router-view class="animated"></router-view>
        </transition>
    </section>
</template>
<script>
    export default {
        name: ''
    }
</script>
<style lang="">

    .app-main {
        height: 100%;
    }
    .show-left-nav .app-main {
        padding-left: 50px;
    }
</style>