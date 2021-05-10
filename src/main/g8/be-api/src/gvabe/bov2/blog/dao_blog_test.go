package blog

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"main/src/gvabe/bov2/user"
	"main/src/utils"
)

const numSampleRows = 100

func initSampleRowsComment(t *testing.T, testName string, dao BlogCommentDao) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numSampleRows; i++ {
		istr := fmt.Sprintf("%03d", i)
		_tagVersion := uint64(1337)
		_id := istr
		_userId := "admin@local"
		_userMaskId := "admin"
		_user := user.NewUser(_tagVersion, _userId, _userMaskId)
		_postIsPublic := true
		_postTitle := "Blog post title"
		_postContent := "Blog post content"
		_post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
		_commentContent := "Blog comment content"
		c := NewBlogComment(_tagVersion, _user, _post, nil, _commentContent)
		c.SetId(_id)
		_numLikes := float64(123)
		c.SetDataAttr("props.tag", "1357")
		c.SetDataAttr("props.active", true)
		c.SetDataAttr("num_likes", _numLikes)
		if ok, err := dao.Create(c); err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", testName+"/Create", ok, err)
		}
	}
}

func doTestCommentDaoCreateGet(t *testing.T, name string, dao BlogCommentDao) {
	_tagVersion := uint64(1337)
	_id := utils.UniqueId()
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)
	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	_post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
	_commentContent := "Blog comment content"
	comment0 := NewBlogComment(_tagVersion, _user, _post, nil, _commentContent)
	comment0.SetId(_id)
	_numLikes := float64(123)
	comment0.SetDataAttr("props.tag", "1357")
	comment0.SetDataAttr("props.active", true)
	comment0.SetDataAttr("num_likes", _numLikes)

	if ok, err := dao.Create(comment0); err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	if comment1, err := dao.Get(_id); err != nil || comment1 == nil {
		t.Fatalf("%s failed: nil or error %s", name+"/Get("+_id+")", err)
	} else {
		if v1, v0 := comment1.GetDataAttrAsUnsafe("props.tag", reddo.TypeString), "1357"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetDataAttrAsUnsafe("props.active", reddo.TypeBool), true; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetDataAttrAsUnsafe("num_likes", reddo.TypeInt), int64(_numLikes); v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetTagVersion(), _tagVersion; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetId(), _id; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetContent(), _commentContent; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetOwnerId(), _userId; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetPostId(), _post.GetId(); v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetParentId(), ""; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if comment1.GetChecksum() != comment0.GetChecksum() {
			t.Fatalf("%s failed: expected %#v but received %#v", name, comment0.GetChecksum(), comment1.GetChecksum())
		}
	}
}

func doTestCommentDaoCreateUpdateGet(t *testing.T, name string, dao BlogCommentDao) {
	_tagVersion := uint64(1337)
	_id := utils.UniqueId()
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)
	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	_post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
	_commentContent := "Blog comment content"
	comment0 := NewBlogComment(_tagVersion, _user, _post, nil, _commentContent)
	comment0.SetId(_id)
	_numLikes := float64(123)
	comment0.SetDataAttr("props.tag", "1357")
	comment0.SetDataAttr("props.active", true)
	comment0.SetDataAttr("num_likes", _numLikes)

	if ok, err := dao.Create(comment0); err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	comment0.SetContent(_commentContent + "-new").SetOwnerId(_userId + "-new").SetPostId(_post.GetId() + "-new").SetParentId("-new").SetTagVersion(_tagVersion + 3)
	comment0.SetDataAttr("props.tag", "2468")
	comment0.SetDataAttr("props.active", false)
	comment0.SetDataAttr("num_likes", _numLikes+2)
	if ok, err := dao.Update(comment0); err != nil {
		t.Fatalf("%s failed: %s", name+"/Update", err)
	} else if !ok {
		t.Fatalf("%s failed: cannot update record", name)
	}
	if comment1, err := dao.Get(_id); err != nil || comment1 == nil {
		t.Fatalf("%s failed: nil or error %s", name+"/Get("+_id+")", err)
	} else {
		if v1, v0 := comment1.GetDataAttrAsUnsafe("props.tag", reddo.TypeString), "2468"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetDataAttrAsUnsafe("props.active", reddo.TypeBool), false; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetDataAttrAsUnsafe("num_likes", reddo.TypeInt), int64(_numLikes+2); v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetTagVersion(), _tagVersion+3; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetId(), _id; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetContent(), _commentContent+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetOwnerId(), _userId+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetPostId(), _post.GetId()+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := comment1.GetParentId(), "-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if comment1.GetChecksum() != comment0.GetChecksum() {
			t.Fatalf("%s failed: expected %#v but received %#v", name, comment0.GetChecksum(), comment1.GetChecksum())
		}
	}
}

func doTestCommentDaoCreateDelete(t *testing.T, name string, dao BlogCommentDao) {
	_tagVersion := uint64(1337)
	_id := utils.UniqueId()
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)
	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	_post := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
	_commentContent := "Blog comment content"
	comment0 := NewBlogComment(_tagVersion, _user, _post, nil, _commentContent)
	comment0.SetId(_id)
	_numLikes := float64(123)
	comment0.SetDataAttr("props.tag", "1357")
	comment0.SetDataAttr("props.active", true)
	comment0.SetDataAttr("num_likes", _numLikes)

	if ok, err := dao.Create(comment0); err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	if comment1, err := dao.Get(_id); err != nil || comment1 == nil {
		t.Fatalf("%s failed: nil or error %s", name+"/Get("+_id+")", err)
	} else if ok, err := dao.Delete(comment1); !ok || err != nil {
		t.Fatalf("%s failed: not-ok or error %s", name+"/Delete("+_id+")", err)
	}

	if comment1, err := dao.Get(_id); err != nil || comment1 != nil {
		t.Fatalf("%s failed: not-nil or error %s", name+"/Get("+_id+")", err)
	}
}

func doTestCommentDaoGetAll(t *testing.T, name string, dao BlogCommentDao) {
	initSampleRowsComment(t, name, dao)
	commentList, err := dao.GetAll(nil, nil)
	if err != nil || len(commentList) != numSampleRows {
		t.Fatalf("%s failed: expected %#v but received %#v (error %s)", name+"/GetAll", numSampleRows, len(commentList), err)
	}
}

func doTestCommentDaoGetN(t *testing.T, name string, dao BlogCommentDao) {
	initSampleRowsComment(t, name, dao)
	commentList, err := dao.GetN(3, 5, nil, nil)
	if err != nil || len(commentList) != 5 {
		t.Fatalf("%s failed: expected %#v but received %#v (error %s)", name+"/GetN", 5, len(commentList), err)
	}
}

/*----------------------------------------------------------------------*/

var userList []*user.User
var userPostCount map[string]int
var userFeedCount map[string]int

func initSampleRowsPost(t *testing.T, testName string, dao BlogPostDao) {
	rand.Seed(time.Now().UnixNano())
	userList = make([]*user.User, 0)
	userPostCount = make(map[string]int)
	userFeedCount = make(map[string]int)
	for i := 0; i < 4; i++ {
		_tagVersion := uint64(1337)
		_userId := strconv.Itoa(i)
		_userMaskId := _userId
		_user := user.NewUser(_tagVersion, _userId, _userMaskId)
		userList = append(userList, _user)
		userPostCount[_userId] = 0
		userFeedCount[_userId] = 0
	}
	for i := 0; i < numSampleRows; i++ {
		istr := fmt.Sprintf("%03d", i)
		_tagVersion := uint64(1337)
		_id := istr
		_user := userList[rand.Intn(len(userList))]
		userPostCount[_user.GetId()]++
		_postIsPublic := rand.Intn(1024)%2 == 0
		if _postIsPublic {
			for k, _ := range userFeedCount {
				userFeedCount[k]++
			}
		} else {
			userFeedCount[_user.GetId()]++
		}
		_postTitle := "Blog post title"
		_postContent := "Blog post content"
		p := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
		p.SetId(_id)
		_numLikes := float64(123)
		p.SetDataAttr("props.tag", "1357")
		p.SetDataAttr("props.active", true)
		p.SetDataAttr("num_likes", _numLikes)
		if ok, err := dao.Create(p); err != nil || !ok {
			t.Fatalf("%s failed: %#v / %s", testName+"/Create", ok, err)
		}
	}
}

func doTestPostDaoCreateGet(t *testing.T, name string, dao BlogPostDao) {
	_tagVersion := uint64(1337)
	_id := utils.UniqueId()
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)
	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	post0 := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
	post0.SetId(_id)
	_numComments := 12
	_numVotesUp := 34
	_numVotesDown := 56
	post0.SetNumComments(_numComments).SetNumVotesUp(_numVotesUp).SetNumVotesDown(_numVotesDown)
	_numLikes := float64(123)
	post0.SetDataAttr("props.tag", "1357")
	post0.SetDataAttr("props.active", true)
	post0.SetDataAttr("num_likes", _numLikes)

	if ok, err := dao.Create(post0); err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	if post1, err := dao.Get(_id); err != nil || post1 == nil {
		t.Fatalf("%s failed: nil or error %s", name+"/Get("+_id+")", err)
	} else {
		if v1, v0 := post1.GetDataAttrAsUnsafe("props.tag", reddo.TypeString), "1357"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetDataAttrAsUnsafe("props.active", reddo.TypeBool), true; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetDataAttrAsUnsafe("num_likes", reddo.TypeInt), int64(_numLikes); v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetTagVersion(), _tagVersion; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetId(), _id; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetContent(), _postContent; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetTitle(), _postTitle; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetOwnerId(), _userId; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetNumComments(), _numComments; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetNumVotesUp(), _numVotesUp; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetNumVotesDown(), _numVotesDown; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if post1.GetChecksum() != post0.GetChecksum() {
			t.Fatalf("%s failed: expected %#v but received %#v", name, post0.GetChecksum(), post1.GetChecksum())
		}
	}
}

func doTestPostDaoCreateUpdateGet(t *testing.T, name string, dao BlogPostDao) {
	_tagVersion := uint64(1337)
	_id := utils.UniqueId()
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)
	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	post0 := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
	post0.SetId(_id)
	_numComments := 12
	_numVotesUp := 34
	_numVotesDown := 56
	post0.SetNumComments(_numComments).SetNumVotesUp(_numVotesUp).SetNumVotesDown(_numVotesDown)
	_numLikes := float64(123)
	post0.SetDataAttr("props.tag", "1357")
	post0.SetDataAttr("props.active", true)
	post0.SetDataAttr("num_likes", _numLikes)

	if ok, err := dao.Create(post0); err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	post0.SetNumComments(_numComments + 1).SetNumVotesUp(_numVotesUp + 2).SetNumVotesDown(_numVotesDown + 3).
		SetTitle(_postTitle + "-new").SetContent(_postContent + "-new").SetPublic(!_postIsPublic).SetOwnerId(_userId + "-new").
		SetTagVersion(_tagVersion + 3)
	post0.SetDataAttr("props.tag", "2468")
	post0.SetDataAttr("props.active", false)
	post0.SetDataAttr("num_likes", _numLikes+2)
	if ok, err := dao.Update(post0); err != nil {
		t.Fatalf("%s failed: %s", name+"/Update", err)
	} else if !ok {
		t.Fatalf("%s failed: cannot update record", name)
	}
	if post1, err := dao.Get(_id); err != nil || post1 == nil {
		t.Fatalf("%s failed: nil or error %s", name+"/Get("+_id+")", err)
	} else {
		if v1, v0 := post1.GetDataAttrAsUnsafe("props.tag", reddo.TypeString), "2468"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetDataAttrAsUnsafe("props.active", reddo.TypeBool), false; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetDataAttrAsUnsafe("num_likes", reddo.TypeInt), int64(_numLikes+2); v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetTagVersion(), _tagVersion+3; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetId(), _id; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetId(), _id; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetContent(), _postContent+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetTitle(), _postTitle+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetOwnerId(), _userId+"-new"; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetNumComments(), _numComments+1; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetNumVotesUp(), _numVotesUp+2; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if v1, v0 := post1.GetNumVotesDown(), _numVotesDown+3; v1 != v0 {
			t.Fatalf("%s failed: expected %#v but received %#v", name, v0, v1)
		}
		if post1.GetChecksum() != post0.GetChecksum() {
			t.Fatalf("%s failed: expected %#v but received %#v", name, post0.GetChecksum(), post1.GetChecksum())
		}
	}
}

func doTestPostDaoCreateDelete(t *testing.T, name string, dao BlogPostDao) {
	_tagVersion := uint64(1337)
	_id := utils.UniqueId()
	_userId := "admin@local"
	_userMaskId := "admin"
	_user := user.NewUser(_tagVersion, _userId, _userMaskId)
	_postIsPublic := true
	_postTitle := "Blog post title"
	_postContent := "Blog post content"
	post0 := NewBlogPost(_tagVersion, _user, _postIsPublic, _postTitle, _postContent)
	post0.SetId(_id)
	_numComments := 12
	_numVotesUp := 34
	_numVotesDown := 56
	post0.SetNumComments(_numComments).SetNumVotesUp(_numVotesUp).SetNumVotesDown(_numVotesDown)
	_numLikes := float64(123)
	post0.SetDataAttr("props.tag", "1357")
	post0.SetDataAttr("props.active", true)
	post0.SetDataAttr("num_likes", _numLikes)

	if ok, err := dao.Create(post0); err != nil || !ok {
		t.Fatalf("%s failed: %#v / %s", name+"/Create", ok, err)
	}

	if post1, err := dao.Get(_id); err != nil || post1 == nil {
		t.Fatalf("%s failed: nil or error %s", name+"/Get("+_id+")", err)
	} else if ok, err := dao.Delete(post1); !ok || err != nil {
		t.Fatalf("%s failed: not-ok or error %s", name+"/Delete("+_id+")", err)
	}

	if post1, err := dao.Get(_id); err != nil || post1 != nil {
		t.Fatalf("%s failed: not-nil or error %s", name+"/Get("+_id+")", err)
	}
}

func doTestPostDaoGetUserPostsAll(t *testing.T, name string, dao BlogPostDao) {
	initSampleRowsPost(t, name, dao)
	postList, err := dao.GetUserPostsAll(userList[0])
	if err != nil || len(postList) != userPostCount[userList[0].GetId()] {
		t.Fatalf("%s failed: expected %#v but received %#v (error %s)", name+"/GetUserPostsAll", userPostCount[userList[0].GetId()], len(postList), err)
	}
}

func doTestPostDaoGetUserPostsN(t *testing.T, name string, dao BlogPostDao) {
	initSampleRowsPost(t, name, dao)
	postList, err := dao.GetUserPostsN(userList[0], 1, 2)
	numExpected := userPostCount[userList[0].GetId()]
	if numExpected < 1 {
		numExpected = 0
	} else if numExpected < 2 {
		numExpected = 1
	} else {
		numExpected = 2
	}
	if err != nil || len(postList) != numExpected {
		t.Fatalf("%s failed: expected %#v but received %#v (error %s)", name+"/GetUserPostsN", numExpected, len(postList), err)
	}
}

func doTestPostDaoGetUserFeedAll(t *testing.T, name string, dao BlogPostDao) {
	initSampleRowsPost(t, name, dao)
	postList, err := dao.GetUserFeedAll(userList[0])
	if err != nil || len(postList) != userFeedCount[userList[0].GetId()] {
		t.Fatalf("%s failed: expected %#v but received %#v (error %s)", name+"/GetUserFeedAll", userFeedCount[userList[0].GetId()], len(postList), err)
	}
}

func doTestPostDaoGetUserFeedN(t *testing.T, name string, dao BlogPostDao) {
	initSampleRowsPost(t, name, dao)
	postList, err := dao.GetUserFeedN(userList[0], 1, 2)
	numExpected := userFeedCount[userList[0].GetId()]
	if numExpected < 1 {
		numExpected = 0
	} else if numExpected < 2 {
		numExpected = 1
	} else {
		numExpected = 2
	}
	if err != nil || len(postList) != numExpected {
		t.Fatalf("%s failed: expected %#v but received %#v (error %s)", name+"/GetUserFeedN", numExpected, len(postList), err)
	}
}
