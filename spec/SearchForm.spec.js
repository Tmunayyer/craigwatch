import { mount, shallowMount, createLocalVue } from '@vue/test-utils';
import VueRouter from "vue-router";
import SearchForm from '../src/components/SearchForm.vue';
import App from '@/App.vue';
import mocker from './mocker.js';
import { routes } from '../src/router.js';

describe("input fields", () => {
    const wrapper = shallowMount(SearchForm);
    // should render top level div
    expect(wrapper.get(".search-form"));

    it("name-field: should alter state", () => {
        const nameInputWrapper = wrapper.get("#name");
        expect(nameInputWrapper);

        const setTo = "bladerunner";
        nameInputWrapper.setValue(setTo);
        expect(wrapper.vm.$data.name).toBe(setTo);
    });

    it("url-field: should alter state", () => {
        const urlInputWrapper = wrapper.get("#url");
        expect(urlInputWrapper);

        const setTo = "https://newyork.craigslist.org/";
        urlInputWrapper.setValue(setTo);
        expect(wrapper.vm.$data.url).toBe(setTo);
    });
});

describe("submit and request", () => {
    const localVue = createLocalVue();
    localVue.use(VueRouter);
    const router = new VueRouter({ routes });

    it("should redirect to results page", async () => {
        const httpReturnData = { ID: 99 };
        const wrapper = shallowMount(SearchForm, {
            router,
            mocks: {
                $http: mocker.api({
                    shouldFail: false,
                    data: httpReturnData
                })
            }
        });

        const nameInputWrapper = wrapper.get("#name");
        const urlInputWrapper = wrapper.get("#url");

        nameInputWrapper.setValue("bladerunner");
        urlInputWrapper.setValue("https://newyork.craigslist.org/");

        wrapper.get(".submit-button").trigger("click");
        await wrapper.vm.$nextTick();

        expect(wrapper.vm.$route.path).toBe(`/result/${httpReturnData.ID}`);
    });
});
