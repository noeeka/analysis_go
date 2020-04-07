<template>
    <div id="app" class="wrapper">
        <Header></Header>
        <Navicate></Navicate>
        <HelloWorld></HelloWorld>

    </div>
</template>

<script>
    import "../plugins/fontawesome-free/css/all.min.css"
    import "../dist/css/adminlte.min.css"
    import Header from "@/components/Header";
    import Navicate from "@/components/Navicate";
    import HelloWorld from './components/Analysis.vue'
    import "../plugins/jquery/jquery.min.js"
    import "../plugins/bootstrap/js/bootstrap.bundle.min.js"
    import "../dist/js/adminlte.js"
    import "../dist/js/demo.js"


    export default {
        name: 'App',
        components: {
            Navicate,
            HelloWorld
        }
    }
</script>

<style>

</style>
