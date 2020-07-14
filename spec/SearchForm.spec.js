import { shallowMount } from '@vue/test-utils';
import SearchForm from '../src/components/SearchForm.vue';

const wrapper = shallowMount(SearchForm);

test("first spec", () => {
    console.log("the wrapper:", wrapper);
});