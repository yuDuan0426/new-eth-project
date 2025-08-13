// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * 合约模板集合
 * 包含常用的智能合约模板，可以直接复制使用
 */

// ============================================================================
// 1. 基础存储合约模板
// ============================================================================

contract BasicStorage {
    uint256 private storedData;
    address public owner;
    
    event DataStored(uint256 newValue, address indexed by);
    
    modifier onlyOwner() {
        require(msg.sender == owner, "Not the owner");
        _;
    }
    
    constructor() {
        owner = msg.sender;
    }
    
    function set(uint256 x) public onlyOwner {
        storedData = x;
        emit DataStored(x, msg.sender);
    }
    
    function get() public view returns (uint256) {
        return storedData;
    }
}

// ============================================================================
// 2. 投票合约模板
// ============================================================================

contract Voting {
    struct Proposal {
        string name;
        uint256 voteCount;
        bool exists;
    }
    
    struct Voter {
        bool hasVoted;
        uint256 proposalId;
        bool isRegistered;
    }
    
    address public chairperson;
    mapping(address => Voter) public voters;
    mapping(uint256 => Proposal) public proposals;
    uint256 public proposalCount;
    bool public votingOpen;
    
    event VoterRegistered(address indexed voter);
    event ProposalAdded(uint256 indexed proposalId, string name);
    event VoteCast(address indexed voter, uint256 indexed proposalId);
    event VotingClosed(uint256 winningProposalId);
    
    modifier onlyChairperson() {
        require(msg.sender == chairperson, "Not the chairperson");
        _;
    }
    
    modifier onlyRegisteredVoter() {
        require(voters[msg.sender].isRegistered, "Not a registered voter");
        _;
    }
    
    modifier votingIsOpen() {
        require(votingOpen, "Voting is closed");
        _;
    }
    
    constructor() {
        chairperson = msg.sender;
        votingOpen = true;
    }
    
    function registerVoter(address voter) public onlyChairperson {
        require(!voters[voter].isRegistered, "Voter already registered");
        voters[voter].isRegistered = true;
        emit VoterRegistered(voter);
    }
    
    function addProposal(string memory name) public onlyChairperson {
        proposals[proposalCount] = Proposal({
            name: name,
            voteCount: 0,
            exists: true
        });
        emit ProposalAdded(proposalCount, name);
        proposalCount++;
    }
    
    function vote(uint256 proposalId) public onlyRegisteredVoter votingIsOpen {
        require(!voters[msg.sender].hasVoted, "Already voted");
        require(proposals[proposalId].exists, "Proposal does not exist");
        
        voters[msg.sender].hasVoted = true;
        voters[msg.sender].proposalId = proposalId;
        proposals[proposalId].voteCount++;
        
        emit VoteCast(msg.sender, proposalId);
    }
    
    function closeVoting() public onlyChairperson {
        votingOpen = false;
        uint256 winningProposalId = getWinningProposal();
        emit VotingClosed(winningProposalId);
    }
    
    function getWinningProposal() public view returns (uint256 winningProposalId) {
        uint256 winningVoteCount = 0;
        for (uint256 i = 0; i < proposalCount; i++) {
            if (proposals[i].voteCount > winningVoteCount) {
                winningVoteCount = proposals[i].voteCount;
                winningProposalId = i;
            }
        }
    }
}

// ============================================================================
// 3. 拍卖合约模板
// ============================================================================

contract Auction {
    address payable public beneficiary;
    uint256 public auctionEndTime;
    
    address public highestBidder;
    uint256 public highestBid;
    
    mapping(address => uint256) public pendingReturns;
    bool public ended;
    
    event HighestBidIncreased(address bidder, uint256 amount);
    event AuctionEnded(address winner, uint256 amount);
    
    modifier onlyBefore(uint256 time) {
        require(block.timestamp < time, "Too late");
        _;
    }
    
    modifier onlyAfter(uint256 time) {
        require(block.timestamp > time, "Too early");
        _;
    }
    
    constructor(uint256 biddingTime, address payable beneficiaryAddress) {
        beneficiary = beneficiaryAddress;
        auctionEndTime = block.timestamp + biddingTime;
    }
    
    function bid() public payable onlyBefore(auctionEndTime) {
        require(msg.value > highestBid, "Bid not high enough");
        
        if (highestBid != 0) {
            pendingReturns[highestBidder] += highestBid;
        }
        
        highestBidder = msg.sender;
        highestBid = msg.value;
        emit HighestBidIncreased(msg.sender, msg.value);
    }
    
    function withdraw() public returns (bool) {
        uint256 amount = pendingReturns[msg.sender];
        if (amount > 0) {
            pendingReturns[msg.sender] = 0;
            
            if (!payable(msg.sender).send(amount)) {
                pendingReturns[msg.sender] = amount;
                return false;
            }
        }
        return true;
    }
    
    function auctionEnd() public onlyAfter(auctionEndTime) {
        require(!ended, "Auction already ended");
        
        ended = true;
        emit AuctionEnded(highestBidder, highestBid);
        
        beneficiary.transfer(highestBid);
    }
}

// ============================================================================
// 4. 众筹合约模板
// ============================================================================

contract Crowdfunding {
    struct Campaign {
        address payable creator;
        uint256 goal;
        uint256 pledged;
        uint256 startAt;
        uint256 endAt;
        bool claimed;
        bool exists;
    }
    
    mapping(uint256 => Campaign) public campaigns;
    mapping(uint256 => mapping(address => uint256)) public pledgedAmount;
    uint256 public campaignCount;
    
    event CampaignCreated(uint256 indexed id, address indexed creator, uint256 goal, uint256 startAt, uint256 endAt);
    event Pledged(uint256 indexed id, address indexed caller, uint256 amount);
    event Unpledged(uint256 indexed id, address indexed caller, uint256 amount);
    event Claimed(uint256 indexed id);
    event Refunded(uint256 indexed id, address indexed caller, uint256 amount);
    
    function createCampaign(uint256 goal, uint256 startAt, uint256 endAt) external {
        require(startAt >= block.timestamp, "Start time is in the past");
        require(endAt >= startAt, "End time is before start time");
        require(endAt <= block.timestamp + 90 days, "End time is too far in the future");
        
        campaigns[campaignCount] = Campaign({
            creator: payable(msg.sender),
            goal: goal,
            pledged: 0,
            startAt: startAt,
            endAt: endAt,
            claimed: false,
            exists: true
        });
        
        emit CampaignCreated(campaignCount, msg.sender, goal, startAt, endAt);
        campaignCount++;
    }
    
    function pledge(uint256 id) external payable {
        Campaign storage campaign = campaigns[id];
        require(campaign.exists, "Campaign does not exist");
        require(block.timestamp >= campaign.startAt, "Campaign has not started");
        require(block.timestamp <= campaign.endAt, "Campaign has ended");
        
        campaign.pledged += msg.value;
        pledgedAmount[id][msg.sender] += msg.value;
        
        emit Pledged(id, msg.sender, msg.value);
    }
    
    function unpledge(uint256 id, uint256 amount) external {
        Campaign storage campaign = campaigns[id];
        require(campaign.exists, "Campaign does not exist");
        require(block.timestamp <= campaign.endAt, "Campaign has ended");
        require(pledgedAmount[id][msg.sender] >= amount, "Insufficient pledged amount");
        
        campaign.pledged -= amount;
        pledgedAmount[id][msg.sender] -= amount;
        
        payable(msg.sender).transfer(amount);
        
        emit Unpledged(id, msg.sender, amount);
    }
    
    function claim(uint256 id) external {
        Campaign storage campaign = campaigns[id];
        require(campaign.exists, "Campaign does not exist");
        require(campaign.creator == msg.sender, "Not the campaign creator");
        require(block.timestamp > campaign.endAt, "Campaign has not ended");
        require(campaign.pledged >= campaign.goal, "Campaign did not reach goal");
        require(!campaign.claimed, "Campaign already claimed");
        
        campaign.claimed = true;
        campaign.creator.transfer(campaign.pledged);
        
        emit Claimed(id);
    }
    
    function refund(uint256 id) external {
        Campaign storage campaign = campaigns[id];
        require(campaign.exists, "Campaign does not exist");
        require(block.timestamp > campaign.endAt, "Campaign has not ended");
        require(campaign.pledged < campaign.goal, "Campaign reached goal");
        
        uint256 bal = pledgedAmount[id][msg.sender];
        pledgedAmount[id][msg.sender] = 0;
        payable(msg.sender).transfer(bal);
        
        emit Refunded(id, msg.sender, bal);
    }
}

// ============================================================================
// 5. 简单 NFT 合约模板
// ============================================================================

interface IERC165 {
    function supportsInterface(bytes4 interfaceId) external view returns (bool);
}

interface IERC721 {
    function balanceOf(address owner) external view returns (uint256 balance);
    function ownerOf(uint256 tokenId) external view returns (address owner);
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes calldata data) external;
    function safeTransferFrom(address from, address to, uint256 tokenId) external;
    function transferFrom(address from, address to, uint256 tokenId) external;
    function approve(address to, uint256 tokenId) external;
    function setApprovalForAll(address operator, bool approved) external;
    function getApproved(uint256 tokenId) external view returns (address operator);
    function isApprovedForAll(address owner, address operator) external view returns (bool);
    
    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
    event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId);
    event ApprovalForAll(address indexed owner, address indexed operator, bool approved);
}

contract SimpleNFT is IERC165, IERC721 {
    string public name;
    string public symbol;
    
    mapping(uint256 => address) private _owners;
    mapping(address => uint256) private _balances;
    mapping(uint256 => address) private _tokenApprovals;
    mapping(address => mapping(address => bool)) private _operatorApprovals;
    
    uint256 private _currentTokenId = 0;
    
    constructor(string memory _name, string memory _symbol) {
        name = _name;
        symbol = _symbol;
    }
    
    function supportsInterface(bytes4 interfaceId) public view virtual override returns (bool) {
        return interfaceId == type(IERC721).interfaceId || interfaceId == type(IERC165).interfaceId;
    }
    
    function balanceOf(address owner) public view virtual override returns (uint256) {
        require(owner != address(0), "ERC721: balance query for the zero address");
        return _balances[owner];
    }
    
    function ownerOf(uint256 tokenId) public view virtual override returns (address) {
        address owner = _owners[tokenId];
        require(owner != address(0), "ERC721: owner query for nonexistent token");
        return owner;
    }
    
    function approve(address to, uint256 tokenId) public virtual override {
        address owner = ownerOf(tokenId);
        require(to != owner, "ERC721: approval to current owner");
        require(msg.sender == owner || isApprovedForAll(owner, msg.sender), "ERC721: approve caller is not owner nor approved for all");
        
        _approve(to, tokenId);
    }
    
    function getApproved(uint256 tokenId) public view virtual override returns (address) {
        require(_exists(tokenId), "ERC721: approved query for nonexistent token");
        return _tokenApprovals[tokenId];
    }
    
    function setApprovalForAll(address operator, bool approved) public virtual override {
        require(operator != msg.sender, "ERC721: approve to caller");
        _operatorApprovals[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }
    
    function isApprovedForAll(address owner, address operator) public view virtual override returns (bool) {
        return _operatorApprovals[owner][operator];
    }
    
    function transferFrom(address from, address to, uint256 tokenId) public virtual override {
        require(_isApprovedOrOwner(msg.sender, tokenId), "ERC721: transfer caller is not owner nor approved");
        _transfer(from, to, tokenId);
    }
    
    function safeTransferFrom(address from, address to, uint256 tokenId) public virtual override {
        safeTransferFrom(from, to, tokenId, "");
    }
    
    function safeTransferFrom(address from, address to, uint256 tokenId, bytes memory _data) public virtual override {
        require(_isApprovedOrOwner(msg.sender, tokenId), "ERC721: transfer caller is not owner nor approved");
        _safeTransfer(from, to, tokenId, _data);
    }
    
    function mint(address to) public returns (uint256) {
        uint256 tokenId = _currentTokenId;
        _currentTokenId++;
        _mint(to, tokenId);
        return tokenId;
    }
    
    function _exists(uint256 tokenId) internal view virtual returns (bool) {
        return _owners[tokenId] != address(0);
    }
    
    function _isApprovedOrOwner(address spender, uint256 tokenId) internal view virtual returns (bool) {
        require(_exists(tokenId), "ERC721: operator query for nonexistent token");
        address owner = ownerOf(tokenId);
        return (spender == owner || getApproved(tokenId) == spender || isApprovedForAll(owner, spender));
    }
    
    function _mint(address to, uint256 tokenId) internal virtual {
        require(to != address(0), "ERC721: mint to the zero address");
        require(!_exists(tokenId), "ERC721: token already minted");
        
        _balances[to] += 1;
        _owners[tokenId] = to;
        
        emit Transfer(address(0), to, tokenId);
    }
    
    function _transfer(address from, address to, uint256 tokenId) internal virtual {
        require(ownerOf(tokenId) == from, "ERC721: transfer from incorrect owner");
        require(to != address(0), "ERC721: transfer to the zero address");
        
        _approve(address(0), tokenId);
        
        _balances[from] -= 1;
        _balances[to] += 1;
        _owners[tokenId] = to;
        
        emit Transfer(from, to, tokenId);
    }
    
    function _approve(address to, uint256 tokenId) internal virtual {
        _tokenApprovals[tokenId] = to;
        emit Approval(ownerOf(tokenId), to, tokenId);
    }
    
    function _safeTransfer(address from, address to, uint256 tokenId, bytes memory _data) internal virtual {
        _transfer(from, to, tokenId);
        require(_checkOnERC721Received(from, to, tokenId, _data), "ERC721: transfer to non ERC721Receiver implementer");
    }
    
    function _checkOnERC721Received(address from, address to, uint256 tokenId, bytes memory _data) private returns (bool) {
        if (to.code.length > 0) {
            try IERC721Receiver(to).onERC721Received(msg.sender, from, tokenId, _data) returns (bytes4 retval) {
                return retval == IERC721Receiver.onERC721Received.selector;
            } catch (bytes memory reason) {
                if (reason.length == 0) {
                    revert("ERC721: transfer to non ERC721Receiver implementer");
                } else {
                    assembly {
                        revert(add(32, reason), mload(reason))
                    }
                }
            }
        } else {
            return true;
        }
    }
}

interface IERC721Receiver {
    function onERC721Received(address operator, address from, uint256 tokenId, bytes calldata data) external returns (bytes4);
}

// ============================================================================
// 6. 代币质押合约模板
// ============================================================================

interface IERC20Simple {
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
}

contract TokenStaking {
    IERC20Simple public stakingToken;
    IERC20Simple public rewardToken;
    
    uint256 public rewardRate = 100; // 每秒奖励代币数量
    uint256 public lastUpdateTime;
    uint256 public rewardPerTokenStored;
    
    mapping(address => uint256) public userRewardPerTokenPaid;
    mapping(address => uint256) public rewards;
    mapping(address => uint256) public balances;
    
    uint256 public totalSupply;
    
    event Staked(address indexed user, uint256 amount);
    event Withdrawn(address indexed user, uint256 amount);
    event RewardPaid(address indexed user, uint256 reward);
    
    constructor(address _stakingToken, address _rewardToken) {
        stakingToken = IERC20Simple(_stakingToken);
        rewardToken = IERC20Simple(_rewardToken);
    }
    
    modifier updateReward(address account) {
        rewardPerTokenStored = rewardPerToken();
        lastUpdateTime = block.timestamp;
        
        if (account != address(0)) {
            rewards[account] = earned(account);
            userRewardPerTokenPaid[account] = rewardPerTokenStored;
        }
        _;
    }
    
    function rewardPerToken() public view returns (uint256) {
        if (totalSupply == 0) {
            return rewardPerTokenStored;
        }
        return rewardPerTokenStored + (((block.timestamp - lastUpdateTime) * rewardRate * 1e18) / totalSupply);
    }
    
    function earned(address account) public view returns (uint256) {
        return ((balances[account] * (rewardPerToken() - userRewardPerTokenPaid[account])) / 1e18) + rewards[account];
    }
    
    function stake(uint256 amount) external updateReward(msg.sender) {
        require(amount > 0, "Cannot stake 0");
        totalSupply += amount;
        balances[msg.sender] += amount;
        stakingToken.transferFrom(msg.sender, address(this), amount);
        emit Staked(msg.sender, amount);
    }
    
    function withdraw(uint256 amount) public updateReward(msg.sender) {
        require(amount > 0, "Cannot withdraw 0");
        require(balances[msg.sender] >= amount, "Insufficient balance");
        totalSupply -= amount;
        balances[msg.sender] -= amount;
        stakingToken.transfer(msg.sender, amount);
        emit Withdrawn(msg.sender, amount);
    }
    
    function getReward() public updateReward(msg.sender) {
        uint256 reward = rewards[msg.sender];
        if (reward > 0) {
            rewards[msg.sender] = 0;
            rewardToken.transfer(msg.sender, reward);
            emit RewardPaid(msg.sender, reward);
        }
    }
    
    function exit() external {
        withdraw(balances[msg.sender]);
        getReward();
    }
}

// ============================================================================
// 7. 简单 DAO 合约模板
// ============================================================================

contract SimpleDAO {
    struct Proposal {
        uint256 id;
        address proposer;
        string description;
        uint256 amount;
        address payable recipient;
        uint256 votes;
        uint256 deadline;
        bool executed;
        mapping(address => bool) voters;
    }
    
    mapping(uint256 => Proposal) public proposals;
    mapping(address => uint256) public shares;
    uint256 public totalShares;
    uint256 public proposalCount;
    uint256 public constant VOTING_PERIOD = 7 days;
    uint256 public constant QUORUM = 51; // 51%
    
    event ProposalCreated(uint256 indexed proposalId, address indexed proposer, string description, uint256 amount, address recipient);
    event Voted(uint256 indexed proposalId, address indexed voter, uint256 shares);
    event ProposalExecuted(uint256 indexed proposalId);
    event SharesIssued(address indexed to, uint256 amount);
    
    modifier onlyMember() {
        require(shares[msg.sender] > 0, "Not a member");
        _;
    }
    
    function issueShares(address to, uint256 amount) external {
        shares[to] += amount;
        totalShares += amount;
        emit SharesIssued(to, amount);
    }
    
    function createProposal(string memory description, uint256 amount, address payable recipient) external onlyMember {
        uint256 proposalId = proposalCount++;
        Proposal storage proposal = proposals[proposalId];
        proposal.id = proposalId;
        proposal.proposer = msg.sender;
        proposal.description = description;
        proposal.amount = amount;
        proposal.recipient = recipient;
        proposal.deadline = block.timestamp + VOTING_PERIOD;
        
        emit ProposalCreated(proposalId, msg.sender, description, amount, recipient);
    }
    
    function vote(uint256 proposalId) external onlyMember {
        Proposal storage proposal = proposals[proposalId];
        require(block.timestamp < proposal.deadline, "Voting period ended");
        require(!proposal.voters[msg.sender], "Already voted");
        
        proposal.voters[msg.sender] = true;
        proposal.votes += shares[msg.sender];
        
        emit Voted(proposalId, msg.sender, shares[msg.sender]);
    }
    
    function executeProposal(uint256 proposalId) external {
        Proposal storage proposal = proposals[proposalId];
        require(block.timestamp >= proposal.deadline, "Voting period not ended");
        require(!proposal.executed, "Proposal already executed");
        require(proposal.votes * 100 >= totalShares * QUORUM, "Quorum not reached");
        require(address(this).balance >= proposal.amount, "Insufficient funds");
        
        proposal.executed = true;
        proposal.recipient.transfer(proposal.amount);
        
        emit ProposalExecuted(proposalId);
    }
    
    receive() external payable {}
}