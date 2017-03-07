import { RepoPage } from './app.po';

describe('repo App', () => {
  let page: RepoPage;

  beforeEach(() => {
    page = new RepoPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
