<template>
    <ol class="breadcrumb breadcrumb-fixed">
        <li class="breadcrumb-item">
            <i class="fa" :class="titleClass"/>&nbsp;{{title}}
        </li>
        <li class="ml-auto">
            <slot/>
        </li>
    </ol>
</template>

<script>
    export default {
        name: "vTitle",
        props: {
            title: String,
            titleClass: {
                type: String,
                default: ""
            },
        }
    }
</script>
