import { mount, shallowMount, createLocalVue } from '@vue/test-utils';
import VueRouter from "vue-router";
import SearchList from '../src/components/SearchList.vue';
import mocker from './mocker.js';
import { routes } from '../src/router.js';



describe("list rendering", () => {
    const localVue = createLocalVue();
    localVue.use(VueRouter);
    const router = new VueRouter({ routes });

    const fakeSearchList = [
        {
            "ID": 22,
            "Name": "Major Tom",
            "URL": "https://newyork.craigslist.org/search/sss?query=ground%20control\u0026sort=rel",
            "CreatedOn": "2020-07-06T18:51:51.516996-04:00",
            "UnixCutoffDate": 1594056360
        }
    ];

    const sl = mount(SearchList, {
        router,
        mocks: {
            $http: new mocker.api({
                shouldFail: false,
                data: fakeSearchList
            })
        }
    });

    it("should populate state on mount", () => {
        expect(sl.vm.$data.searchList.length).toBe(1);
    });

    it("should render the list item", () => {
        const item = sl.get(".search-listitem");
        expect(item);

        const header = item.get(".header-name");
        expect(header.text()).toBe(fakeSearchList[0].Name);
    });

    it("should redirect on click", async () => {
        const item = sl.get(".search-listitem");
        expect(item);

        item.trigger("click");
        await sl.vm.$nextTick();

        expect(sl.vm.$route.path).toBe(`/result/${fakeSearchList[0].ID}`);
    });
});

describe("api error state", () => {
    it("should not display the error", async () => {
        const sl = mount(SearchList, {
            mocks: {
                $http: new mocker.api({
                    shouldFail: false,
                    data: []
                })
            }
        });

        // wait for async request to finish
        await sl.vm.$nextTick();
        expect(sl.vm.$data.error).toBe(false);
        expect(sl.findComponent({ name: "Error" }).exists()).toBe(false);
    });

    it("should display the error", async () => {
        const sl = mount(SearchList, {
            mocks: {
                $http: new mocker.api({
                    shouldFail: true,
                    data: []
                })
            }
        });

        // wait for async request to finish
        await sl.vm.$nextTick();
        expect(sl.vm.$data.error).toBe(true);
        expect(sl.findComponent({ name: "Error" }).exists());
    });

});