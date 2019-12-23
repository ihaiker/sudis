<template>
    <nav v-if="pages.length > 1" aria-label="navigation">
        <ul class="pagination justify-content-center">
            <li class="page-link">
                共 {{items.total}} 条
            </li>

            <li class="page-item" :class="{'disabled':start}">
                <button class="page-link" @click="toPage(1)">首页</button>
            </li>

            <li v-for="(p) in pages" class="page-item" :class="{'active':(p === page)}">
                <button class="page-link" @click="toPage(p)">{{p}}</button>
            </li>

            <li class="page-item" :class="{'disabled':end}">
                <button class="page-link" @click="toPage(pages.length)">尾页</button>
            </li>

            <li class="page-link">
                每页 {{items.limit}} 条
            </li>
        </ul>
    </nav>
</template>
<script>
    export default {
        name: "XPage",
        props: {items: Object},
        data: () => ({
            page: 1, limit: 12,
            start: false, end: false, pages: [],
        }),
        created() {
            this.set()
        },
        methods: {
            toPage(page) {
                this.$emit("change", {page: page, limit: this.limit});
            },
            set(value) {
                try {
                    this.page = value.page;
                    this.limit = value.limit;
                    this.pages = [];
                    let maxPage = Math.ceil(value.total / value.limit);
                    for (let i = 1; i <= maxPage; i++) {
                        this.pages.push(i);
                    }
                    this.start = (value.page === 1);
                    this.end = (value.page === this.pages.length);
                } catch (e) {
                }
            }
        },
        watch: {
            items(value) {
                this.set(value);
            }
        }
    }
</script>
