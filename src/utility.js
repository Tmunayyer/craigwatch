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
    },

    chartDate: function (date) {
        const translations = {
            1: "st",
            2: "nd",
            3: "rd",
        };
        const day = new Date(date).getUTCDate();
        const time = date.slice(11);
        let daySuffix = 'th';
        let remainder = day % 10;
        if (translations[day] !== undefined) {
            // 1, 2, 3
            daySuffix = trnaslations[day];
        } else if (day > 3 && day < 20) {
            // 4 - 19
        } else if (translations[remainder] !== undefined) {
            // 20+
            daySuffix = translations[remainder];
        }
        return `${day}${daySuffix} ${time}`;
    }
}