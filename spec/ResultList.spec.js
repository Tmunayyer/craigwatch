import { mount } from '@vue/test-utils';
import ResultList from '../src/components/ResultList.vue';
import mocker from './mocker.js';

describe("list rendering", () => {
    const fakeListings = {
        "HasNewListings": true,
        "Listings": [
            {
                "ID": 17,
                "SearchID": 22,
                "DataPID": "7149585315",
                "DataRepostOf": "7047865536",
                "UnixDate": 1594056360,
                "Title": "TASCAM DM3200 Digital Mixer w/ Parametric EQ Effects Pro Tools Control",
                "Link": "https://newyork.craigslist.org/wch/msg/d/carmel-tascam-dm3200-digital-mixer/7149585315.html",
                "Price": 750,
                "Hood": " (Westchester Putnam County)"
            },
            {
                "ID": 28,
                "SearchID": 22,
                "DataPID": "7154598663",
                "DataRepostOf": "",
                "UnixDate": 1594055160,
                "Title": "2017 Yamaha SCR 950 *LEFTOVER SALE!*",
                "Link": "https://newyork.craigslist.org/lgi/mcd/d/plainfield-2017-yamaha-scr-950-leftover/7154598663.html",
                "Price": 6495,
                "Hood": " (Motorsports Nation Plainfield)"
            }
        ]
    };

    it("should render the listings", async () => {
        const list = mount(ResultList, {
            mocks: {
                $http: mocker.api({
                    shouldFail: false,
                    data: fakeListings
                })
            }
        });
        // required to await when mounting inside a test
        await list.vm.$nextTick();
        expect(list.exists());
        expect(list.vm.$data.resultList.length).toBe(2);
    });

    it("should render the listings", async () => {
        const list = mount(ResultList, {
            mocks: {
                $http: mocker.api({
                    shouldFail: true,
                    data: []
                })
            }
        });
        // 1st tick sets all the data
        await list.vm.$nextTick();
        expect(list.exists());
        expect(list.vm.$data.error).toBe(true);
        expect(list.vm.$data.resultList.length).toBe(0);

        // 2nd tick renders the components, maybe the subcomponents?
        await list.vm.$nextTick();
        expect(list.findComponent({ name: "Error" }).exists()).toBe(true);
    });

    it("should correctly manage the interval", async () => {
        const list = mount(ResultList, {
            mocks: {
                $http: mocker.api({
                    shouldFail: false,
                    data: fakeListings
                })
            }
        });

        // required to await when mounting inside a test
        await list.vm.$nextTick();
        expect(list.exists());
        expect(list.vm.$data.resultList.length).toBe(2);

        // hold a referene to the polling obj
        const data = list.vm.$data;
        expect(data.polling).not.toBe(null);

        list.destroy();
        expect(list.exists()).toBe(false);
        expect(data.polling).toBe(null);
    });
});
