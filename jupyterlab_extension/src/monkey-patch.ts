import { Launcher } from '@jupyterlab/launcher';
import * as React from 'react';

export function monkeyPatchLauncher() {
  const customCategories = ['Templates'];

  const originalRender = (Launcher as any).prototype.render;

  function customRender(this: Launcher): React.ReactElement<any> | null {
    const originalReactElement = originalRender.call(
      this
    ) as React.ReactElement<any>;

    function moveCustomCategories(
      element: React.ReactElement<any>
    ): React.ReactElement<any> {
      if (element.props && element.props.className === 'jp-Launcher-content') {
        const outerChildren = element.props.children.slice();
        const innerChildren = outerChildren[1].slice();

        const reorderedChildren = customCategories
          .map(customCategory => {
            const index = innerChildren.findIndex(
              (child: React.ReactElement<any>) => child.key === customCategory
            );
            if (index >= 0) {
              return innerChildren.splice(index, 1)[0];
            }
            return null;
          })
          .filter(child => child !== null)
          .concat(innerChildren);

        const newOuterChildren = [outerChildren[0], reorderedChildren];
        return React.cloneElement(element, element.props, newOuterChildren);
      }

      if (element.props && element.props.children) {
        const newChildren = React.Children.map(
          element.props.children,
          (child: React.ReactElement<any>) => {
            return moveCustomCategories(child);
          }
        );
        return React.cloneElement(element, element.props, newChildren);
      }

      return element;
    }

    const modifiedReactElement = moveCustomCategories(originalReactElement);

    return modifiedReactElement;
  }

  (Launcher as any).prototype.render = customRender;
}
