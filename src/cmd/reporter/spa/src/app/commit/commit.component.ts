import { Input, Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-commit',
  templateUrl: './commit.component.html',
  styleUrls: ['./commit.component.less']
})
export class CommitComponent implements OnInit {
 @Input() userId: number;
  @Input() name: string;
  @Input() email: string;
  @Input() message: string;
  @Input() avatar: string;
  @Input() commit: string;
  @Input() date: Date;
  tags: number[] = [1059];
  nameIn = '';
  viewMessage = '';
  isLoaded = false;
  constructor() {

  }
  mock() {
    this.userId = 1;
    this.name = 'Kian Ostad';
    this.email = "c@kianostad.com";
    this.message = `codegen changed for data tables 
    
    fix #234`
    this.avatar = '1ebe42c3d1195c7b588e0abe9c178fa7';
    this.commit = 'c6fe905c47e6ab8454fdb0fcf88417210d71f05e'
    this.date = new Date();
  }
  ngOnInit() {
    this.mock()
    this.viewMessage = this.message.substr(0, this.message.indexOf('\n'))
    this.nameIn = this.name.substr(0, 2)
    const time = Math.random()
    setTimeout(() => this.isLoaded = true, this.randomInt(800,1200))
  }
  getTags() {
    var re = new RegExp("(fix|ref)? +#([0-9]+)", "gmi")
    var m;
    while (m = re.exec(this.message)) {
      console.log(m[1], m[2]);
    }
  }
  randomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1) + min);
  }

}
