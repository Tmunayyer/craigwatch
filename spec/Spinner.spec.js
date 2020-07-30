import { mount } from '@vue/test-utils';
import Spinner, { spinnerState } from '../src/components/Spinner.vue';

describe("spinner state handlers", () => {
    jest.useFakeTimers();

    it("handles loading correctly", async () => {
        const wrapper = mount(Spinner);
        expect(wrapper);

        spinnerState.setLoading(true);
        expect(spinnerState.state.timeout).not.toBe(null);
        expect(spinnerState.state.loading).toBe(false);

        jest.runAllTimers();
        expect(setTimeout).toBeCalled();
        expect(spinnerState.state.loading).toBe(true);


        spinnerState.setLoading(false);
        expect(spinnerState.state.timeout).toBe(null);
        expect(spinnerState.state.loading).toBe(false);
    });
});
