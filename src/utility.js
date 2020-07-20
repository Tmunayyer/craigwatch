export default {
    /**
     * formatDate accepts any form of date and returns
     *  a presentable string
     * 
     * @param {*} date 
     * @returns {string}
     */
    formatDate: function (date) {
        const options = {
            year: "numeric",
            month: "numeric",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit"
        };
        const today = new Date(date);

        return today.toLocaleDateString("en-US", options);

    }
}