import utils from "./utils";

let mixins = {
    data: () => ({
        gid: ('m' + Math.floor(Math.random() * 10000000000)),
        timers: {}
    }),
    methods: {
        gvid(prefix, subfix) {
            if (undefined !== prefix && undefined !== subfix) {
                return prefix + this.gid + subfix;
            } else if (prefix) {
                return this.gid + prefix;
            } else {
                return this.gid;
            }
        },
        loading(isLoading, title) {
            this.$root.loadingShow = isLoading;
            this.$root.loadingTitle = title;
        },
        startLoading(title) {
            this.loading(true, title || "loading......");
        },
        loaddingStatus(title) {
            this.$root.loadingTitle = title;
        },
        finishLoading() {
            this.loading(false, "");
        },
        request(title, request) {
            this.startLoading(title);
            request.finally(this.finishLoading)
        },
        ramShow(limit) {
            return utils.ram(limit)
        },
        now() {
            return utils.now()
        },
    }
};
export default  mixins;
