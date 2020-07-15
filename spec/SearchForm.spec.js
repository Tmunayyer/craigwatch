import { mount, shallowMount, createLocalVue } from '@vue/test-utils';
import VueRouter from "vue-router";
import SearchForm from '../src/components/SearchForm.vue';
import mocker from './mocker.js';
import { routes } from '../src/router.js';

describe("input fields", () => {
    const wrapper = mount(SearchForm);
    // should render top level div
    expect(wrapper.get(".search-form"));

    it("name-field: should alter state", () => {
        const nameInputWrapper = wrapper.get("#name");
        expect(nameInputWrapper);

        const inputField = nameInputWrapper.get("input");
        expect(inputField);

        const setTo = "bladerunner";
        inputField.setValue(setTo);
        expect(wrapper.vm.$data.name).toBe(setTo);
    });

    it("url-field: should alter state", () => {
        const urlInputWrapper = wrapper.get("#url");
        expect(urlInputWrapper);

        const inputField = urlInputWrapper.get("input");
        expect(inputField);

        const setTo = "https://newyork.craigslist.org/";
        inputField.setValue(setTo);
        expect(wrapper.vm.$data.url).toBe(setTo);
    });
});

describe("input field validations", () => {
    // NOTE: testing error boundries on validation means it should never send an api request
    // and never redirect. These should still be mocked in the event of improper handling
    // during a refactor etc.

    const localVue = createLocalVue();
    localVue.use(VueRouter);
    const router = new VueRouter({ routes });

    it("should detect empty fields", async () => {
        const wrapper = mount(SearchForm, {
            mocks: {
                $http: mocker.api({
                    shouldFail: true,
                    data: {},
                })
            }
        });
        expect(wrapper.get(".search-form"));

        wrapper.get("button").trigger("click");
        await wrapper.vm.$nextTick();

        expect(wrapper.vm.$data.nameErr).toBe(true);
        expect(wrapper.vm.$data.urlErr).toBe(true);
    });

    it("should validate the url", async () => {
        const wrapper = mount(SearchForm, {
            mocks: {
                $http: mocker.api({
                    shouldFail: true,
                    data: {},
                })
            }
        });
        expect(wrapper.get(".search-form"));

        wrapper.get("#name").get("input").setValue("valid name");
        wrapper.get("#url").get("input").setValue("invalid url");

        wrapper.get("button").trigger("click");
        await wrapper.vm.$nextTick();

        expect(wrapper.vm.$data.nameErr).toBe(false);
        expect(wrapper.vm.$data.urlErr).toBe(true);
    });
});

describe("submit and request", () => {
    const localVue = createLocalVue();
    localVue.use(VueRouter);
    const router = new VueRouter({ routes });

    it("should redirect to results page", async () => {
        const httpReturnData = { ID: 99 };
        const wrapper = mount(SearchForm, {
            router,
            mocks: {
                $http: mocker.api({
                    shouldFail: false,
                    data: httpReturnData
                })
            }
        });

        const nameInputWrapper = wrapper.get("#name");
        const nameInput = nameInputWrapper.get("input");
        const urlInputWrapper = wrapper.get("#url");
        const urlInput = urlInputWrapper.get("input");

        nameInput.setValue("bladerunner");
        urlInput.setValue("https://newyork.craigslist.org/");

        wrapper.get("button").trigger("click");
        await wrapper.vm.$nextTick();

        expect(wrapper.vm.$route.path).toBe(`/result/${httpReturnData.ID}`);
    });
});
