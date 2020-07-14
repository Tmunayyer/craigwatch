import { shallowMount } from '@vue/test-utils';
import SearchForm from '../src/components/SearchForm.vue';

const wrapper = shallowMount(SearchForm);

describe("input fields", () => {
    // should render top level div
    expect(wrapper.get(".search-form"));

    test("name-field: should alter state", () => {
        const nameInputWrapper = wrapper.get("#name");
        expect(nameInputWrapper);

        const setTo = "bladerunner";
        nameInputWrapper.setValue(setTo);
        expect(wrapper.vm.$data.name).toBe(setTo);
    });

    test("url-field: should alter state", () => {
        const urlInputWrapper = wrapper.get("#name");
        expect(urlInputWrapper);

        const setTo = "https://newyork.craigslist.org/";
        urlInputWrapper.setValue(setTo);
        expect(wrapper.vm.$data.name).toBe(setTo);
    });

    // test("submit: should make request", () => {
    //     const nameInputWrapper = wrapper.get("#name");
    //     const urlInputWrapper = wrapper.get("#name");

    //     nameInputWrapper.setValue("balderunner");
    //     urlInputWrapper.setValue("https://newyork.craigslist.org/");



    // })
});
