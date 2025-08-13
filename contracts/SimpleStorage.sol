// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title SimpleStorage
 * @dev Store & retrieve value in a variable
 */
contract SimpleStorage {
    uint256 storedData;

    /**
     * @dev Store value in variable
     * @param x value to store
     */
    function set(uint256 x) public {
        storedData = x;
    }

    /**
     * @dev Return value 
     * @return value of 'storedData'
     */
    function get() public view returns (uint256) {
        return storedData;
    }

    /**
     * @dev Increment stored value by 1
     */
    function increment() public {
        storedData += 1;
    }

    /**
     * @dev Reset stored value to 0
     */
    function reset() public {
        storedData = 0;
    }
}