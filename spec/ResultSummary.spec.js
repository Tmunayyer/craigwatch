import { mount } from '@vue/test-utils';
import ResultSummary from '../src/components/ResultSummary.vue';
import mocker from './mocker.js';

describe("should display data", () => {
    const fakeSearchDetails = {
        "ID": 22,
        "Name": "Major Tom",
        "URL": "https://newyork.craigslist.org/search/sss?query=ground%20control\u0026sort=rel",
        "CreatedOn": "2020-07-06T18:51:51.516996-04:00",
        "UnixCutoffDate": 1594056360
    };

    it("should render", async () => {
        const summary = mount(ResultSummary, {
            mocks: {
                $http: mocker.api({
                    shouldFail: false,
                    data: fakeSearchDetails
                })
            }
        });

        await summary.vm.$nextTick();
        expect(summary.get(".result-header"));

        await summary.vm.$nextTick();
        expect(summary.vm.$data.searchDetails.Name).toBe(fakeSearchDetails.Name);
        const name = summary.get(".result-header-name");
        expect(name);
        expect(name.text()).toBe(fakeSearchDetails.Name);
    });

    it("should handle api erros", async () => {
        const summary = mount(ResultSummary, {
            mocks: {
                $http: mocker.api({
                    shouldFail: true,
                    data: {}
                })
            }
        });
        await summary.vm.$nextTick();
        expect(summary);

        await summary.vm.$nextTick();
        expect(summary.vm.$data.error).toBe(true);

        expect(summary.findComponent({ name: "Error" }).exists()).toBe(true);
    });
});