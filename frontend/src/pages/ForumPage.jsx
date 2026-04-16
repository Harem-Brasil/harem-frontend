import { useOutletContext } from "react-router-dom";
import Composer from "../components/Composer";
import FeedList from "../components/FeedList";

export default function ForumPage() {
  const { posts, publishPost } = useOutletContext();
  const forumPosts = posts.filter((post) => post.tags.includes("forum"));

  return (
    <>
      <Composer
        title="Criar topico no forum"
        pageType="forum"
        onPublish={publishPost}
      />
      <FeedList posts={forumPosts} />
    </>
  );
}
