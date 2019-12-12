<template>
    <div v-if="isLoadingShow" class="loader">
        <div class="content">
            <div class="dot"></div>
            <div class="dot"></div>
            <div class="dot"></div>
            <div class="dot"></div>
            <div class="dot"></div>
        </div>
        <div class="header">{{$root.loadingTitle}}</div>
    </div>
</template>

<script>
    export default {
        name: "loadings",
        computed: {
            isLoadingShow() {
                return this.$root.loadingShow;
            }
        }
    }
</script>

<style lang="scss">
    .loader {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        z-index: 99999998;

        .content, .header {
            position: absolute;
            top: 50%;
            left: 50%;
            padding: 5rem;
            transform: translate3d(-50%, -50%, 0);
        }

        .header {
            font-size: 1rem;
            font-weight: bold;
        }
    }

    .dot {
        width: 3rem;
        height: 3rem;
        background: #3ac;
        border-radius: 100%;
        display: inline-block;
        animation: slide 1s infinite;

        &:nth-child(1) {
            animation-delay: 0.1s;
            background: #32aacc;
        }

        &:nth-child(2) {
            animation-delay: 0.2s;
            background: #64aacc;
        }

        &:nth-child(3) {
            animation-delay: 0.3s;
            background: #96aacc;
        }

        &:nth-child(4) {
            animation-delay: 0.4s;
            background: #c8aacc;
        }

        &:nth-child(5) {
            animation-delay: 0.5s;
            background: #faaacc;
        }
    }

    @mixin keyframes($animation) {
        @-webkit-keyframes #{$animation} {
            @content;
        }
        @-moz-keyframes #{$animation} {
            @content;
        }
        @-o-keyframes #{$animation} {
            @content;
        }
        @keyframes #{$animation} {
            @content;
        }
    }

    @include keyframes(slide) {
        0% {
            transform: scale(1);
            opacity: 0;
        }
        50% {
            opacity: 0.3;
            transform: scale(2);
        }
        100% {
            transform: scale(1);
            opacity: 1;
        }
    }
</style>
