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

        let adjustedTimestamp = new Date(date);
        // if (typeof date === "number") {
        //     const today = new Date(date);
        //     const minutes = today.getTimezoneOffset();
        //     console.log("the today:", today);
        //     console.log("the offset:", minutes);
        //     const seconds = minutes * 60;
        //     const tzOffset = seconds * 1000;

        //     adjustedTimestamp = new Date(date + tzOffset);
        // }

        const output = adjustedTimestamp.toLocaleString("en-US");
        return output;
    },

    /**
     * Uses unix timestamp only.
     * 
     * @param {interger} date 
     */
    chartDate: function (date) {
        const options = {
            year: "numeric",
            month: "numeric",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit"
        };
        const translations = {
            1: "st",
            2: "nd",
            3: "rd",
        };

        // const today = new Date(date);
        // const minutes = today.getTimezoneOffset();
        // const seconds = minutes * 60;
        // const tzOffset = seconds * 1000;

        const adjustedTimestamp = new Date(date);
        const localString = adjustedTimestamp.toLocaleDateString("en-US", options);
        // console.log("the localstring:", localString);
        const day = adjustedTimestamp.getDate();

        const time = localString.slice(11);
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